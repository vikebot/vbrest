package vbcore

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// CryptoGenThreshold defines the maximum number N to retry generating
// cryptographically secure random bytes inside CryptoGenBytes, if previous
// ones fail
var CryptoGenThreshold = 10

// CryptoGen generates a random 32 bytes long AES key and return's it's base64
// equivalent
func CryptoGen() (key string, err error) {
	buf, err := CryptoGenBytes(32)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

// CryptoGenBytes uses CSPRNGs or HSMs to generate `len` cryptographically
// random bytes that can be used in any security specific context. In the
// rare event that no secure random values could be generated this func
// will retry it N-times (`n` defined as `CryptoGenThreshold` - defaults to
// `10`). If after N tries no crypo randoms are found the error will be
// returned and the slice will be empty.
func CryptoGenBytes(len int) (buf []byte, err error) {
	if CryptoGenThreshold < 1 {
		return []byte{}, fmt.Errorf("vbcore.CryptoGenBytes: CryptoGenThreshold mustn't be less than 1 but is %v", CryptoGenThreshold)
	}

	buf = make([]byte, len)
	for i := 0; i < CryptoGenThreshold; i++ {
		_, err = io.ReadFull(rand.Reader, buf)
		if err == nil {
			return
		}
	}
	return []byte{}, err
}

const cryptoGenStringSet = "abcdefghijklmnopqrstuvwxyzABZDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// CryptoGenString generates a n-long random string that can be used in any
// security specific context. Internally CryptoGenBytes it used (read func-
// docs for more info).
func CryptoGenString(n int) (string, error) {
	buf, err := CryptoGenBytes(n)
	if err != nil {
		return "", err
	}

	str := ""
	for _, b := range buf {
		str += string(cryptoGenStringSet[(b>>2)%62])
	}

	return str, nil
}
