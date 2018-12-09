package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	logSimple "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/valyala/fasthttp"
	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbjwt"
	"github.com/vikebot/vbnet"
	"github.com/vikebot/vbrest/vbapi"
	"github.com/vikebot/vbrest/vbmail"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	stat statsd.Statter
)

func main() {
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	// get config
	if configFlag == nil || len(*configFlag) == 0 {
		logSimple.Fatalln("argument '-config' is mandatory and mustn't be empty")
	}
	configBuf, err := ioutil.ReadFile(*configFlag)
	if err != nil {
		logSimple.Fatalln(err)
	}
	config := &conf{}
	err = json.Unmarshal(configBuf, config)
	if err != nil {
		logSimple.Fatalln(err)
	}

	// Logging server
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	console := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	logCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, console, priority),
	)
	log = zap.New(logCore)

	// Statsd aggregation funcs
	log.Info("init statsd")
	stat, err = statsd.NewBufferedClient("", "", time.Second*1, 0)
	if err != nil {
		log.Error("unable to connect to statsd server. creating noop statter")
		stat, _ = statsd.NewNoopClient()
	}

	// Init our database connection
	log.Info("init vbapi")
	err = vbapi.Init(config.Recaptcha.Secret, config.DB.Addr, config.DB.User, config.DB.Pass, config.DB.Name, log)
	if err != nil {
		log.Fatal("unable to init db connection", zap.Error(err))
	}

	// Init our sendgrid client
	log.Info("init vbmail")
	vbmail.Init(config.Sendgrid.Secret)

	// Print CORS message
	if config.CORS.Enabled {
		if config.CORS.Wildcard {
			log.Warn("cors enabled with wildcard")
		} else {
			log.Info("cors enabled with set of domains", zap.Strings("allowed", config.CORS.AllowedDomains))
		}
	}

	// Fill our rt (routes-tree) -> Radix tree has better lookup-times
	// than iterating over each specified route in the hashmap (direct
	// hashmap-lookups cannot be performed (we need prefix-matching!)
	log.Info("init endpoints")
	rt := iradix.New()
	for _, ep := range allEndpoints(log) {
		var failed bool
		rt, _, failed = rt.Insert([]byte(ep.Name), ep)
		if failed {
			log.Fatal("unable to insert route into routes-tree", zap.String("route", ep.Name))
		}
	}

	// Load all our signing keys used for validating the JWTs sent from
	// clients to authenticate themselves.
	log.Info("init vbjwt")
	vbjwt.Init(config.JWT.ProductionIsssuer, config.JWT.DefaultSigningKeyID, config.JWT.SigningKeys, log)

	respond := func(c *fasthttp.RequestCtx, r interface{}, ctx *zap.Logger) {
		// If r == nil we where succesful so set response: ok
		if r == nil {
			r = &simpleResponse{Response: "ok"}
		}

		// Check for environment variable to enable local development
		if config.CORS.Enabled {
			origin := string(c.Request.Header.Peek("Origin"))

			var allowed bool

			// allow wildcard access
			if config.CORS.Wildcard {
				allowed = true
			} else {
				for _, domain := range config.CORS.AllowedDomains {
					if origin == domain {
						allowed = true
					}
				}
			}

			if allowed {
				c.Response.Header.Add("Access-Control-Allow-Origin", origin)
				c.Response.Header.Add("Access-Control-Allow-Credentials", "true")
			}
		}

		switch v := r.(type) {
		// Valid request - response only needs to be marshaled and sent
		default:
			body, err := json.Marshal(v)
			if err != nil {
				ctx.Error("marshaling response failed", zap.Error(err))
				stat.Inc("vbrest.response_marshal_error", 1, 1)
				c.SetStatusCode(fasthttp.StatusInternalServerError)
				fmt.Fprint(c, `{"error":"Internal server error"}`)
				return
			}
			resp := string(body)
			ctx.Debug("req_response", zap.String("resp", resp))
			stat.Inc("vbrest.req_ok", 1, 1)
			c.SetStatusCode(fasthttp.StatusOK)
			fmt.Fprint(c, resp)
			return
			// Valid request - but internal server error
		case error:
			if http, ok := v.(vbnet.HTTPError); ok {
				ctx.Info("req_failed", zap.Error(http))
				stat.Inc("vbrest.req_failed_http"+strconv.Itoa(http.HTTPCode()), 1, 1)
				stat.Inc("vbrest.req_failed_code"+strconv.Itoa(http.Code()), 1, 1)
				c.SetStatusCode(http.HTTPCode())
				fmt.Fprintf(c, `{"error":"%s"}`, http.Message())
				return
			}
			ctx.Error("internal_error", zap.Error(v))
			stat.Inc("vbrest.internal_error", 1, 1)
			c.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprint(c, `{"error":"Internal server error"}`)
			return
		}
	}

	dispatch := func(c *fasthttp.RequestCtx) {
		rqid := vbcore.FastRandomString(32)
		ctx := log.With(zap.String("rqid", rqid))

		// Catch panics during execution
		defer func() {
			if rval := recover(); rval != nil {
				ctx.Error("recover handler panic",
					zap.Stack("rval_stack"),
					zap.String("rval_string", fmt.Sprint(rval)))

				c.SetStatusCode(http.StatusInternalServerError)
			}
		}()

		// Convert request buffer to a readable URL-string
		p := string(c.Path())

		// Log basic request meta informations
		ctx.Info("req",
			zap.String("ip", realipFromFasthttp(c)),
			zap.String("path", p))

		// Log the request
		stat.Inc("vbrest.req", 1, 1)

		// Set the response type to json
		c.SetContentType("application/json")

		// Find route
		k, f, ok := rt.Root().LongestPrefix(c.Path())

		// No route matches the request
		if !ok {
			respond(c, errUnknownEndpoit, ctx)
			return
		}

		// Assert reqHandler
		ep, ok := f.(endpoint)
		if !ok {
			log.Error("endpoint assertion failed", zap.String("route", string(k)))
			respond(c, errEndpointAssertionFailed, ctx)
			return
		}

		// If the request handler wants an exact match fail if it would only
		// be a prefix
		if ep.ExactMatch && ep.Name != p {
			respond(c, errUnknownEndpoit, ctx)
			return
		} else if ep.Handler == nil {
			respond(c, errNotImplemented, ctx)
			return
		}

		// https://stackoverflow.com/a/21783145/6123704
		method := string(c.Method())
		if method == "OPTIONS" {
			c.Response.Header.Add("Access-Control-Allow-Credentials", "true")
			c.Response.Header.Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Response.Header.Add("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type, Authorization")
			c.Response.Header.Add("Access-Control-Max-Age", "86400")
			respond(c, nil, ctx)
			return
		}

		// Execute request
		r, err := ep.Handler(c, p, ctx)
		if err != nil {
			respond(c, err, ctx)
			return
		}
		respond(c, r, ctx)
	}

	if config.TLS.Active {
		log.Info("rest service started with https ...",
			zap.String("addr", config.Addr),
			zap.String("cert", config.TLS.Cert),
			zap.String("key", config.TLS.Key))
		err = fasthttp.ListenAndServeTLS(config.Addr, config.TLS.Cert, config.TLS.Key, dispatch)
		if err != nil {
			log.Fatal("ListenAndServeTLS failed", zap.Error(err))
		}
	} else {
		log.Info("rest service started with http ...", zap.String("addr", config.Addr))
		err = fasthttp.ListenAndServe(config.Addr, dispatch)
		if err != nil {
			log.Fatal("ListenAndServe failed", zap.Error(err))
		}
	}
}

func realipFromFasthttp(c *fasthttp.RequestCtx) string {
	// Get ip from x-forwarded-for header
	x := string(c.Request.Header.Peek("X-FORWARDED-FOR"))

	// If we have multiple ips seperated by ',' only use the first one
	if strings.Contains(x, ",") {
		return x[0:strings.Index(x, ",")]
	}

	// Otherwise we only have one -> return it
	return x
}
