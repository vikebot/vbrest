package vbjwt

import (
	"strconv"
	"time"

	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"go.uber.org/zap"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// GenerateNonDefaultCtx creates a new signed token and saves it's JTI into
// the database
func GenerateNonDefaultCtx(issuer string, userID int, expires time.Time, ip string, allowedIPs []string, ctx *zap.Logger) (token string, success bool) {
	// Creation time
	issuedAt := time.Now()

	// Create custom claim
	claims := &VBClaims{
		AllowedIPs: allowedIPs,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   strconv.Itoa(userID),
			Audience:  "api.vikebot.com",
			ExpiresAt: expires.Unix(),
			IssuedAt:  issuedAt.Unix(),
			Id:        vbcore.FastRandomString(32),
		},
	}

	// Safe JTI so we can blacklist it laters
	success = vbdb.JwtAddCtx(claims.Id, expires, userID, issuedAt, ip, ctx)
	if !success {
		return "", false
	}

	// Create JWT and set our current primary signingkey id to the
	// vikebot-signingkey-id in the JWT header.
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	t.Header["vbskid"] = skid

	// Cryptographically sign the token
	st, err := t.SignedString(skstore[skid])
	if err != nil {
		// Log everything we know about the event
		ctx.Error("unable to sign token",
			zap.Int("user_id", userID),
			zap.String("ip", ip),
			zap.Strings("allowed_ips", allowedIPs),
			zap.String("skid", skid),
			zap.Error(err))
		return "", false
	}

	return st, true
}

// GenerateNonDefault is the same as `GenerateNonDefaultCtx` but uses the
// `defaultCtx` as logger.
func GenerateNonDefault(issuer string, userID int, expires time.Time, ip string, allowedIPs []string) (token string, success bool) {
	return GenerateNonDefaultCtx(issuer, userID, expires, ip, allowedIPs, defaultCtx)
}

// GenerateCtx creates a new signed token and saves it's JTI into the database
func GenerateCtx(userID int, ip string, allowedIPs []string, ctx *zap.Logger) (token string, success bool) {
	return GenerateNonDefaultCtx(defaultIssuer, userID, time.Now().Add(time.Hour*24*31), ip, allowedIPs, ctx)
}

// Generate is the same as `GenerateCtx` but uses the `defaultCtx` as logger.
func Generate(userID int, ip string, allowedIPs []string) (token string, success bool) {
	return GenerateCtx(userID, ip, allowedIPs, defaultCtx)
}
