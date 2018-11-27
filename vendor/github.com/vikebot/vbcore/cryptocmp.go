package vbcore

import (
	"crypto/sha512"
	"crypto/subtle"
)

func cryptoCmp(x, y []byte) bool {
	return subtle.ConstantTimeCompare(x, y) == 1
}

// CryptoCmpStr performs a constant-time-compare on strings to prevent from
// side-channel attacks (especiall timing attacks).
// If the strings are not of the same length, they are passed through SHA512
// to normalize their length.
func CryptoCmpStr(x, y string) bool {
	return CryptoCmpBytes([]byte(x), []byte(y))
}

// CryptoCmpBytes performs a constant-time-compare on byte slices to prevent
// from side-channel attacks (especiall timing attacks).
// If the slices are not of the same length, they are passed through SHA512
// to normalize their length.
func CryptoCmpBytes(x, y []byte) bool {
	if len(x) == len(y) {
		return cryptoCmp(x, y)
	}

	// Normalize lengths of x and y and compare their hashes
	xh := sha512.Sum512(x)
	yh := sha512.Sum512(y)
	return cryptoCmp(xh[:], yh[:])
}
