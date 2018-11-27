package main

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbrest/vbapi"
	"go.uber.org/zap"
)

func v1Test(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	return &simpleResponse{Response: "ok"}, nil
}

func v1UserGet(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	userID, err := authproxy(req, vbcore.PermissionDefault, ctx)
	if err != nil {
		return nil, err
	}
	return vbapi.UserGet(userID, ctx)
}

func v1UserGetPublicByID(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	return vbapi.UserGetPublicByID(p[len("/v1/user/get/id/"):], ctx)
}

func v1UserGetPublicByUsername(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	return vbapi.UserGetPublicByUsername(p[len("/v1/user/get/username/"):], ctx)
}

func v1UserUpdate(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	userID, err := authproxy(req, vbcore.PermissionDefault, ctx)
	if err != nil {
		return nil, err
	}

	var user vbcore.User
	err = json.Unmarshal(req.PostBody(), &user)
	if err != nil {
		return nil, err
	}

	err = vbapi.UserUpdate(userID, &user, "", ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func v1RoundActive(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	return vbapi.RoundActive(ctx)
}

func v1RoundJoin(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	userID, err := authproxy(req, vbcore.PermissionDefault, ctx)
	if err != nil {
		return nil, err
	}
	err = vbapi.RoundJoin(userID, p[len("/v1/round/join/"):], ctx)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func v1RoundentryActive(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	userID, err := authproxy(req, vbcore.PermissionDefault, ctx)
	if err != nil {
		return nil, err
	}

	roundentries, err := vbapi.RoundentryActive(userID, ctx)
	if err != nil {
		return nil, err
	}

	return roundentries, nil
}

func v1RoundentryConnectinfo(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	return vbapi.RoundentryConnectinfo(p[len("/v1/roundentry/connectinfo/"):], ctx)
}

func v1RoundentryWatchresolve(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	websocket, err := vbapi.RoundentryWatchresolve(p[len("/v1/roundentry/watchresolve/"):], ctx)
	if err != nil {
		return nil, err
	}
	return websocket, nil
}

func v1RegisterConfirm(req *fasthttp.RequestCtx, p string, ctx *zap.Logger) (r interface{}, err error) {
	ctx.Debug("", zap.String("json", string(req.PostBody())))

	var data vbapi.RegisterConfirmRequest
	err = json.Unmarshal(req.PostBody(), &data)
	if err != nil {
		return nil, err
	}

	err = vbapi.RegisterConfirm(data, realipFromFasthttp(req), ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
