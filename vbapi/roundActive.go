package vbapi

import (
	"errors"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
)

// RoundActive returns a list of all rounds that are currently active
func RoundActive(ctx *zap.Logger) (r []vbcore.Round, err error) {
	var success bool
	if ctx == nil {
		r, success = vbdb.ActiveRounds()
	} else {
		r, success = vbdb.ActiveRoundsCtx(ctx)
	}
	if !success {
		return nil, errors.New("Internal server error")
	}
	return r, nil
}
