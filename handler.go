package main

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type handler func(ctx *fasthttp.RequestCtx, p string, logCtx *zap.Logger) (r interface{}, err error)
