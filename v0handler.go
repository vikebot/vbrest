package main

import (
	"encoding/base32"
	"errors"
	"strconv"

	"github.com/valyala/fasthttp"
	"github.com/vikebot/vbjwt"
	"go.uber.org/zap"
)

var (
	genjwtsecret string
)

func v0AdminGenjwtkey(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	prefixLen := len("/v0/admin/genjwtkey/")
	secretLen := base32.StdEncoding.EncodedLen(128)
	totLen := prefixLen + secretLen + len("/")

	if len(p) < totLen {
		return &simpleResponse{Response: "error: invalid format"}, nil
	}

	// check secret
	secret := p[prefixLen:][:secretLen]
	if secret != genjwtsecret {
		return &simpleResponse{Response: "error: forbidden. unauthorized"}, nil
	}

	// extract userID
	userID, err := strconv.Atoi(p[totLen:])
	if err != nil {
		return &simpleResponse{Response: "error: invalid user_id"}, nil
	}

	// generate token! only for localhost
	token, success := vbjwt.GenerateCtx(userID, "127.0.0.1", []string{"127.0.0.1"}, ctx)
	if !success {
		return nil, errors.New("error during token generation")
	}
	return &simpleResponse{Response: token}, nil
}
