package vbjwt

import (
	"encoding/hex"

	"go.uber.org/zap"
)

var (
	defaultCtx *zap.Logger

	isProd        bool
	defaultIssuer string

	skid    string
	skstore map[string][]byte
)

// Init prepares `vbjwt` for handling authentication later. `isProduction`
// indicates if the default `Generate` methods use "vikebot_production" or
// "vikebot_qa" as JWT issuer. The `signingKeys` map contains all ever used
// keys as [id]hexkeyformat. If you want to deprecate a single key-id remove
// the hexkey and use a empty string.
func Init(isProduction bool, defaultSigningkey string, signingKeys map[string]string, ctx *zap.Logger) error {
	// Save production mode flag and set defaulIssuer accourdingly
	isProd = isProduction
	if isProduction {
		defaultIssuer = "vikebot_production"
	} else {
		defaultIssuer = "vikebot_debug"
	}

	// Store default signing-key and init general store
	skid = defaultSigningkey
	skstore = make(map[string][]byte)
	for k, v := range signingKeys {
		// Check if the key is marked as deprecated
		if len(v) == 0 {
			skstore[k] = nil
		}

		// Try to decode string
		buffer, err := hex.DecodeString(v)
		if err != nil {
			return err
		}

		// Store keyid and buffer in our memory store
		skstore[k] = buffer
	}

	// Save default logging context
	defaultCtx = ctx

	return nil
}
