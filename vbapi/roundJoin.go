package vbapi

import (
	"net/http"
	"strconv"

	"github.com/vikebot/vbdb"
	"github.com/vikebot/vbnet"
	"go.uber.org/zap"
)

// RoundJoin adds the user to the specified round
func RoundJoin(userID int, roundID string, ctx *zap.Logger) error {
	// Validate userID
	if userID < 1 {
		return vbnet.NewHTTPError("User id must be greater than 0", http.StatusBadRequest, codeInvalidUserIDFormat, nil)
	}

	// Validate roundID
	round, err := strconv.Atoi(roundID)
	if err != nil {
		return vbnet.NewHTTPError("Round id must be an uint", http.StatusBadRequest, codeInvalidRoundIDFormat, nil)
	}
	if round < 1 {
		return vbnet.NewHTTPError("Round id must be greater than 0", http.StatusBadRequest, codeInvalidRoundIDFormat, nil)
	}

	exists, success := vbdb.RoundExistsCtx(round, ctx)
	if !success {
		return errInternalServerError
	}
	if !exists {
		return vbnet.NewHTTPError("Specified round doesn't exist", http.StatusBadRequest, codeRoundNotExists, nil)
	}

	// Join
	alreadyJoined, success := vbdb.JoinRoundCtx(userID, round, ctx)
	if !success {
		return errInternalServerError
	}
	if alreadyJoined {
		return vbnet.NewHTTPError("User already joined this round", http.StatusForbidden, codeAlreadyJoined, nil)
	}

	return nil
}
