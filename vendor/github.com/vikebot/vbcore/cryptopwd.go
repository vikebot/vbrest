package vbcore

import (
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    = 4
	argon2Memory  = 8 * 1024
	argon2Threads = 4
	argon2Keylen  = 32
)

func cryptoPwd(password string, salt []byte) []byte {
	sha := sha512.Sum512([]byte(password))
	return argon2.Key(sha[:], salt, argon2Time, argon2Memory, argon2Threads, argon2Keylen)
}

// CryptoPwd performs length normalization and password hashing to the passed
// string and returns the final hash and salt encoded in a hex string
func CryptoPwd(password string) (hash, salt string, err error) {
	saltBuf, err := CryptoGenBytes(32)
	if err != nil {
		return "", "", err
	}

	hashBuf := cryptoPwd(password, saltBuf)
	return hex.EncodeToString(hashBuf), hex.EncodeToString(saltBuf), nil
}

// CryptoPwdVerify compares the already hashed password and it's salt against
// the unhashed given password. This function is resitent from side-channel
// attacks and uses constant-time-comparison methods.
func CryptoPwdVerify(hash, salt, password string) (ok bool, err error) {
	hashBuf, err := hex.DecodeString(hash)
	if err != nil {
		return false, nil
	}

	saltBuf, err := hex.DecodeString(salt)
	if err != nil {
		return false, err
	}

	verifyBuf := cryptoPwd(password, saltBuf)
	return CryptoCmpBytes(hashBuf, verifyBuf), nil
}
