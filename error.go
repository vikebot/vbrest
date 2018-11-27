package main

import (
	"github.com/valyala/fasthttp"
	"github.com/vikebot/vbnet"
)

const (
	codeInternalServerError     = 9000
	codeNotFound                = 9001
	codeNotImplemented          = 9002
	codeInsufficientPermission  = 9003
	codeEndpointAssertionFailed = 9004
)

var (
	errInternalServerError = vbnet.NewHTTPError(
		"Internal Server Error",
		fasthttp.StatusInternalServerError,
		codeInternalServerError,
		nil)
	errUnknownEndpoit = vbnet.NewHTTPError(
		"No API endpoint matches your request",
		fasthttp.StatusNotFound,
		codeNotFound,
		nil)
	errNotImplemented = vbnet.NewHTTPError(
		"Not implemented",
		fasthttp.StatusNotImplemented,
		codeNotImplemented,
		nil)
	errEndpointAssertionFailed = vbnet.NewHTTPError(
		"Internal Server Error",
		fasthttp.StatusInternalServerError,
		codeEndpointAssertionFailed,
		nil)
)
