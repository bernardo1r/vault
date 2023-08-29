package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

const gcmNonceLen = 12

func NewGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	return gcm, err
}

func GCMEncrypt(key []byte, nonce []byte, plaintext []byte) ([]byte, error) {
	gcm, err := NewGCM(key)
	if err != nil {
		return nil, err
	}

	if nonce == nil {
		nonce = make([]byte, gcmNonceLen)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

func GCMDecrypt(key []byte, nonce []byte, ciphertext []byte) ([]byte, error) {
	gcm, err := NewGCM(key)
	if err != nil {
		return nil, err
	}

	if nonce == nil {
		nonce = make([]byte, gcmNonceLen)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}
