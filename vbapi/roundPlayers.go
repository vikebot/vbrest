package vbapi

import (
	"errors"
	"net/http"

	"github.com/vikebot/vbnet"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
)

// RoundPlayers returns a list of all players in a given round
func RoundPlayers(watchtoken string, ctx *zap.Logger) (p []vbcore.Player, err error) {
	var success bool

	if !watchtokenValidator.MatchString(watchtoken) {
		return nil, vbnet.NewHTTPError("Invalid watchtoken format", http.StatusBadRequest, codeInvalidWatchtokenFormat, nil)
	}

	if ctx == nil {
		p, success = vbdb.RoundPlayersWatchtoken(watchtoken)
	} else {
		p, success = vbdb.RoundPlayersWatchtokenCtx(watchtoken, ctx)
	}
	if !success {
		return nil, errors.New("Internal server error")
	}
	return p, nil
}
