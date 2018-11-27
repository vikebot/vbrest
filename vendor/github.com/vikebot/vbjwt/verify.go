package vbjwt

import (
	"strconv"

	"github.com/vikebot/vbdb"
	"github.com/vikebot/vbnet"
	"go.uber.org/zap"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// VerifyCtx authenticates the validity of an JWT token against a set of
// predefined rules. If error interface may be of type `vjwt/(*JWTError)`.
// Only if there aren't any issues `userID` and `permission` will be filled.
// Don't forget to check for `vbcore.PermissionBanned` before using.
func VerifyCtx(token string, ip string, ctx *zap.Logger) (userID int, permission int, err error) {
	// Check if the token is empty
	if len(token) == 0 {
		return 0, 0, errJwtEmpty
	}

	// Parse the token with our custom VBClaims
	jwtToken, err := jwt.ParseWithClaims(token, &VBClaims{}, func(token *jwt.Token) (interface{}, error) {
		var ok bool

		// Check if the token used the correct SigingMethod. In our case HS512
		// (HMAC with SHA2-512)
		if _, ok = token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errJwtUnexpectedSigningMethod
		}

		// See if the token has a vikebot-signing-key-id
		var kid string
		if kid, ok = token.Header["vbskid"].(string); !ok {
			return nil, errJwtTooOld
		}

		// Check if the signing key exists and if it isn't deprecated
		var key []byte
		if key, ok = skstore[kid]; !ok {
			return nil, errJwtTooOld
		}
		if len(key) == 0 {
			return nil, errJwtTooOld
		}

		// No problems -> return key for signing check
		return key, nil
	})
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			// JWT Malformed / Expired / unsinged Signature
			if e.Errors&jwt.ValidationErrorMalformed != 0 {
				return 0, 0, errJwtMalformed
			} else if e.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, 0, errJwtExpired
			} else if e.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return 0, 0, errJwtInvalidSignature
			}
		}

		// Check for custom errors happend during key-fetch
		if e, ok := err.(vbnet.HTTPError); ok {
			return 0, 0, e
		}

		// If error is unknown log it and return unverifiable
		ctx.Warn("unhandled jwt parsing error. Returning unverifiable", zap.Error(err))
		return 0, 0, errJwtUnverifiable
	}

	// Check whether the token is valid or not
	if !jwtToken.Valid {
		return 0, 0, errJwtInvalid
	}

	// As the token is valid we should be able to type assert it to a VBClaims
	// instance
	claims, ok := jwtToken.Claims.(*VBClaims)
	if !ok {
		ctx.Error("unable to type-assert claims to (*VBClaims)")
		return 0, 0, errInternalServerError
	}

	// Check if the issuer of the JWT is the production env
	if claims.Issuer != defaultIssuer {
		ctx.Warn("invalid issuer",
			zap.String("issuer_got", claims.Issuer),
			zap.String("issuer_want", defaultIssuer))
		return 0, 0, errJwtInvalidIssuer
	}

	// Check if the audience of the JWT is the api service
	if claims.Audience != "api.vikebot.com" {
		ctx.Warn("invalid audience",
			zap.String("audience_got", claims.Audience),
			zap.String("audience_want", "api.vikebot.com"))
		return 0, 0, errJwtInvalidAudience
	}

	// Check if the request IP matches the expected origin IP. Should further
	// extend the protections againt JWT theft
	ipAllowed := false
	for _, allowed := range claims.AllowedIPs {
		if allowed == "*" {
			ipAllowed = true
			break
		} else if allowed == ip {
			ipAllowed = true
			break
		}
	}
	if !ipAllowed {
		ctx.Warn("unauthorized origin used JWT",
			zap.String("ip", ip),
			zap.Strings("allowed_ips", claims.AllowedIPs))
		return 0, 0, errUnauthorizedRequestOrigin
	}

	// Convert the JWT subject to our userID
	userID, err = strconv.Atoi(claims.Subject)
	if err != nil {
		ctx.Error("invalid userID in JWT claim",
			zap.String("subject", claims.Subject),
			zap.Error(err))
		return 0, 0, errInternalServerError
	}

	// Check if the JWT is blacklisted
	blacklisted, success := vbdb.JwtIsBlacklistedCtx(claims.Id, ctx)
	if blacklisted {
		ctx.Warn("request with blacklisted jwt",
			zap.String("jti", claims.Id),
			zap.Int("user_id", userID))
		return 0, 0, errJwtBlacklisted
	}

	// Load the user's permission from the database
	userPerm, success := vbdb.UserPermissionCtx(userID, ctx)
	if !success {
		return 0, 0, errInternalServerError
	}

	return userID, userPerm, nil
}

// Verify is the same as `VerifyCtx` but uses the `defaultCtx` as logger.
func Verify(token, string, ip string) (userID int, permission int, err error) {
	return VerifyCtx(token, ip, defaultCtx)
}
