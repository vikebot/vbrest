package vbapi

import (
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
)

// Init configures everything needed for vbapi to work
func Init(recaptchaSecret, dbAddr, dbUser, dbPass, dbName string, ctx *zap.Logger) error {
	recaptcha.Init(recaptchaSecret)

	return vbdb.Init(&vbdb.Config{
		DbAddr: vbcore.NewEndpointAddr(dbAddr),
		DbUser: dbUser,
		DbPass: dbPass,
		DbName: dbName,
	}, ctx)
}
