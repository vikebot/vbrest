package vbapi

import (
	"errors"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
)

// UserUpdate updates an old user profile in vbdb with a new user
// profile
func UserUpdate(userID int, newUser *vbcore.User, msg string, ctx *zap.Logger) error {
	oldUser, success := vbdb.UserFromID(userID)
	if !success {
		return errors.New("Internal server error")
	}

	success = vbdb.UpdateUserCtx(newUser, oldUser, msg, ctx)
	if !success {
		return errors.New("Internal server error")
	}

	return nil
}
