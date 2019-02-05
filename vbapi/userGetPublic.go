package vbapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"github.com/vikebot/vbnet"
	"go.uber.org/zap"
)

// UserGetPublicByID returns the public user profile from the database
// associated with the `userID`.
func UserGetPublicByID(userID string, ctx *zap.Logger) (user *vbcore.SafeUser, err error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, vbnet.NewHTTPError("user_id must be an int", http.StatusBadRequest, codeUserIDMustBeInt, nil)
	}

	user, success := vbdb.UserFromIDCtx(id, ctx)
	if !success {
		return nil, errInternalServerError
	}
	if user == nil {
		return nil, vbnet.NewHTTPError("user_id doesn't exist", http.StatusNotFound, codeUserIDDoesnotExist, nil)
	}

	// Remove any sensitive data from the user
	user.MakePublic()

	return user, nil
}

// UserGetInfoByWatchtoken returns a vbcore.PlayerInfo to display a player in
// the frontend
func UserGetInfoByWatchtoken(watchtoken string, ctx *zap.Logger) (user *vbcore.PlayerInfo, err error) {
	var success bool

	if !watchtokenValidator.MatchString(watchtoken) {
		return nil, vbnet.NewHTTPError("Invalid watchtoken format", http.StatusBadRequest, codeInvalidWatchtokenFormat, nil)
	}

	if ctx == nil {
		user, success = vbdb.UserFromWatchtoken(watchtoken)
	} else {
		user, success = vbdb.UserFromWatchtokenCtx(watchtoken, ctx)
	}
	if !success {
		return nil, errors.New("Internal server error")
	}
	return user, nil
}

// UserGetPublicByUsername returns the public user profile from the database
// associated with the `username`.
func UserGetPublicByUsername(username string, ctx *zap.Logger) (user *vbcore.SafeUser, err error) {
	user, success := vbdb.UserFromUsernameCtx(username, ctx)
	if !success {
		return nil, errInternalServerError
	}
	if user == nil {
		return nil, vbnet.NewHTTPError("username doesn't exist", http.StatusOK, codeUsernameDoesnotExist, nil)
	}

	// Remove any sensitive data from the user
	user.MakePublic()

	return user, nil
}
