package vbapi

import (
	"net/http"
	"regexp"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"github.com/vikebot/vbnet"
	"go.uber.org/zap"
)

var (
	authtokenValidator = regexp.MustCompile("^[a-zA-z0-9]{18}$")
)

// RoundentryConnectinfo trades an authtoken for a
// `vbcore.RoundentryConnectinfo`
func RoundentryConnectinfo(authtoken string, ctx *zap.Logger) (response *vbcore.RoundentryConnectinfo, err error) {
	if !authtokenValidator.MatchString(authtoken) {
		return nil, vbnet.NewHTTPError("Invalid authtoken format", http.StatusBadRequest, codeInvalidAuthtokenFormat, nil)
	}

	connectinfo, exists, success := vbdb.RoundentryConnectinfoCtx(authtoken, ctx)
	if !success {
		return nil, errInternalServerError
	}
	if !exists {
		return nil, vbnet.NewHTTPError("Unknown authtoken", http.StatusOK, codeUnknownAuthtoken, nil)
	}

	return connectinfo, nil
}
