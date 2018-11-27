package vbapi

import (
	"errors"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
)

// UserGet returns the full user profile from the database
func UserGet(userID int, ctx *zap.Logger) (*vbcore.SafeUser, error) {
	user, success := vbdb.UserFromIDCtx(userID, ctx)
	if !success || user == nil {
		return nil, errors.New("Internal server error")
	}

	return user, nil
}
