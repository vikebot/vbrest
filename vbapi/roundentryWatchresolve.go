package vbapi

import (
	"net/http"
	"regexp"

	"github.com/vikebot/vbdb"
	"github.com/vikebot/vbnet"
	"go.uber.org/zap"
)

var (
	watchtokenValidator = regexp.MustCompile("^[a-zA-z0-9]{12}$")
)

// RoundentryWatchresolve resolves a watchtoken to it's associated websocket
// address, so clients only knowing the watchtoken can find their gameserver
// information publisher
func RoundentryWatchresolve(watchtoken string, ctx *zap.Logger) (websocket string, err error) {
	if !watchtokenValidator.MatchString(watchtoken) {
		return "", vbnet.NewHTTPError("Invalid watchtoken format", http.StatusBadRequest, codeInvalidWatchtokenFormat, nil)
	}

	websocket, exists, success := vbdb.WebsocketAddressFromWatchtokenCtx(watchtoken, ctx)
	if !success {
		return "", errInternalServerError
	}
	if !exists {
		return "", vbnet.NewHTTPError("Unknown watchtoken.", http.StatusForbidden, codeUnknownWatchtoken, nil)
	}

	return websocket, nil
}
