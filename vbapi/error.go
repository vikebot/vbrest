package vbapi

import (
	"net/http"

	"github.com/vikebot/vbnet"
)

const (
	codeInternalServerError     = 11000
	codeUserIDMustBeInt         = 11001
	codeUserIDDoesnotExist      = 11002
	codeUsernameDoesnotExist    = 11003
	codeInvalidAuthtokenFormat  = 11004
	codeUnknownAuthtoken        = 11005
	codeInvalidUserIDFormat     = 11006
	codeInvalidRoundIDFormat    = 11007
	codeAlreadyJoined           = 11008
	codeRoundNotExists          = 11009
	codeInvalidWatchtokenFormat = 11010
	codeUnknownWatchtoken       = 11011

	codeInvalidRegisterCode           = 11010
	codeUserCannotBeNull              = 11011
	codeBadUserState                  = 11012
	codeAlreadyFinishedRegistration   = 11013
	codeRegistrationCodeUnknown       = 11014
	codeEmailAddressManipulated       = 11015
	codeEmailStatusManipulated        = 11016
	codeCannotUseMultiplePrimaryEmail = 11017
	codeEmailQuotaExhausted           = 11018
	codeInvalidEmailVerificationCode  = 11019
	codeRegisterVerificationEntry     = 11020
	codeMustHavePrimaryEmail          = 11021
	codeManipulatedWebLink            = 11022
	codeInvalidSocialPlatfrom         = 11023
	codeRecaptchaNotTicked            = 11024
)

var (
	errInternalServerError = vbnet.NewHTTPError(
		"Internal Server Error",
		http.StatusInternalServerError,
		codeInternalServerError,
		nil)
)
