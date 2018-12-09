package main

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbjwt"
	"github.com/vikebot/vbnet"
	"go.uber.org/zap"
)

func authproxy(req *fasthttp.RequestCtx, minPermission int, ctx *zap.Logger) (userID int, err error) {
	var token string

	// Authentication via bearer header
	if header := string(req.Request.Header.Peek("authorization")); len(header) > 0 {
		if !strings.HasPrefix(header, "bearer ") {
			if split := strings.Split(header, " "); len(split) == 2 {
				token = split[1]
			}
		}
		// Authentication via vbauth cookie
	} else if t := string(req.Request.Header.Cookie("vbauth")); len(t) > 0 {
		token = t
	} else {
		return 0, vbnet.NewHTTPError("No auth provided. Access forbidden", fasthttp.StatusForbidden, codeNoAuthProvided, nil)
	}

	userID, permission, err := vbjwt.VerifyCtx(token, realipFromFasthttp(req), ctx)
	if err != nil {
		return 0, err
	}
	ctx.Info("authorized",
		zap.Int("user_id", userID),
		zap.Int("permission", permission))

	// Assert user permission
	if permission < minPermission {
		ctx.Warn("insufficient permission", zap.Int("permission_want", minPermission))

		return 0, vbnet.NewHTTPError(
			fmt.Sprintf("Insufficient permission. Needed %v, has %v", vbcore.PermissionItoA(minPermission), vbcore.PermissionItoA(permission)),
			fasthttp.StatusForbidden,
			codeInsufficientPermission,
			nil)
	}

	return userID, nil
}
