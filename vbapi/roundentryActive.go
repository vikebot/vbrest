package vbapi

import (
	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
)

// RoundentryActive loads all roundentry infos for a given user
func RoundentryActive(userID int, ctx *zap.Logger) (response []vbcore.Roundentry, err error) {
	roundentries, success := vbdb.ActiveRoundentriesCtx(userID, ctx)
	if !success {
		return nil, errInternalServerError
	}

	return roundentries, nil
}
