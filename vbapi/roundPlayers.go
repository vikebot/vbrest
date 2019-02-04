package vbapi

import (
	"errors"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
)

// RoundPlayers returns a list of all players in a given round
func RoundPlayers(roundID int, ctx *zap.Logger) (p []vbcore.Player, err error) {
	var success bool
	if ctx == nil {
		p, success = vbdb.RoundPlayers(roundID)
	} else {
		p, success = vbdb.RoundPlayersCtx(roundID, ctx)
	}
	if !success {
		return nil, errors.New("Internal server error")
	}
	return p, nil
}
