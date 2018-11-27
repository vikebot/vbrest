package vbjwt

import (
	"net/http"

	"github.com/vikebot/vbnet"
)

var (
	errJwtEmpty                   = vbnet.NewHTTPError("Empty JWTs aren't allowed", http.StatusBadRequest, 10000, nil)
	errJwtTooOld                  = vbnet.NewHTTPError("JWT signing key too old. Please refresh your token.", http.StatusBadRequest, 10001, nil)
	errJwtMalformed               = vbnet.NewHTTPError("JWT malformed", http.StatusBadRequest, 10002, nil)
	errJwtExpired                 = vbnet.NewHTTPError("JWT already expired", http.StatusForbidden, 10003, nil)
	errJwtInvalidSignature        = vbnet.NewHTTPError("JWT signature invalid", http.StatusForbidden, 10004, nil)
	errJwtUnverifiable            = vbnet.NewHTTPError("JWT unverifiable", http.StatusBadRequest, 10005, nil)
	errJwtInvalid                 = vbnet.NewHTTPError("JWT token invalid", http.StatusForbidden, 10006, nil)
	errJwtInvalidAudience         = vbnet.NewHTTPError("JWT isn't for this service", http.StatusForbidden, 10007, nil)
	errJwtInvalidIssuer           = vbnet.NewHTTPError("JWT is from an untrusted issuer", http.StatusForbidden, 10008, nil)
	errJwtBlacklisted             = vbnet.NewHTTPError("JWT already blacklisted", http.StatusForbidden, 10009, nil)
	errInternalServerError        = vbnet.NewHTTPError("Internal server error", http.StatusInternalServerError, 10010, nil)
	errUnauthorizedRequestOrigin  = vbnet.NewHTTPError("Unauthorized request origin. Your IP isn't allowed to use this JWT", http.StatusForbidden, 10011, nil)
	errJwtUnexpectedSigningMethod = vbnet.NewHTTPError("Unexpected signing method. Want HS512", http.StatusBadRequest, 10013, nil)
)
