package vbjwt

import "gopkg.in/dgrijalva/jwt-go.v3"

// VBClaims is vikebot's custom `Claims` interface, containing allowed origin
// (e.g. remot) IPs allowed to use this JWT. UserID is stored as `int` in the
// Subject. The `jti` (JWT-ID) can be used to determine blacklisted tokens.
type VBClaims struct {
	AllowedIPs []string `json:"allowed_ips"`
	jwt.StandardClaims
}
