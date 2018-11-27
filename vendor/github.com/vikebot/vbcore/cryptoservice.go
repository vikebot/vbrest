package vbcore

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// CryptoService provides a clean interface to encrypt data
type CryptoService struct {
	gcm cipher.AEAD
}

// NewCryptoService generates a new `CryptoService` instance used to perform
// cryptographically secure encryption and decryption processes using `AES256`
// in `GCM` (Galois-Counter-Mode).
func NewCryptoService(key []byte) (cs *CryptoService, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &CryptoService{
		gcm: gcm,
	}, nil
}

// Encrypt takes the plain content, encrypts it and returns the cipher already
// including the used iv
func (cs CryptoService) Encrypt(plain []byte) (cipher []byte, err error) {
	if len(plain) == 0 {
		return nil, errors.New("vbcore/CryptoService: no content to encrypt")
	}

	nonce, err := CryptoGenBytes(cs.gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	cipher = cs.gcm.Seal(nil, nonce, plain, nil)
	cipher = append(nonce, cipher...)

	return cipher, nil
}

// EncryptBase64 works like `Encrypt`, but encodes the cipher to it's `base64`
// equivalent before returning.
func (cs CryptoService) EncryptBase64(plain []byte) (base64Cipher []byte, err error) {
	cipher, err := cs.Encrypt(plain)
	if err != nil {
		return nil, err
	}

	base64Cipher = make([]byte, base64.RawStdEncoding.EncodedLen(len(cipher)))
	base64.RawStdEncoding.Encode(base64Cipher, cipher)

	return
}

// Decrypt takes the cipher content, decrypts it and returns the plain text
func (cs CryptoService) Decrypt(cipher []byte) (plain []byte, err error) {
	if len(cipher) == 0 {
		return nil, errors.New("vbcore/CryptoService: no content to decrypt")
	}

	nonce := cipher[0:cs.gcm.NonceSize()]
	ciphertext := cipher[cs.gcm.NonceSize():]

	plain, err = cs.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return
}

// DecryptBase64 works like `Decrypt`, but decodes the cipher to it's non
// `base64` equivalent before decrypting.
func (cs CryptoService) DecryptBase64(base64Cipher []byte) (plain []byte, err error) {
	cipher := make([]byte, base64.RawStdEncoding.DecodedLen(len(base64Cipher)))
	_, err = base64.RawStdEncoding.Decode(cipher, base64Cipher)
	if err != nil {
		return nil, err
	}

	return cs.Decrypt(cipher)
}
