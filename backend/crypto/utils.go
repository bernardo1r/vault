package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/sha3"
)

const passwordLen = 32

func DecodePassword(password string) ([]byte, error) {
	passwordDecoded, err := base64.StdEncoding.DecodeString(password)
	switch {
	case err != nil:
		return nil, err

	case len(passwordDecoded) != passwordLen:
		return nil, errors.New("wrong password length")
	}
	return passwordDecoded, nil
}

func CheckPublicKey(publicKey string) error {
	pubkey, _ := base64.StdEncoding.DecodeString(publicKey)
	_, err := x509.ParsePKIXPublicKey(pubkey)
	return err
}

func VerifySignature(publicKeySign string, data []byte, sig []byte) error {
	pubkeyDecoded, _ := base64.StdEncoding.DecodeString(publicKeySign)
	pubkey, err := x509.ParsePKIXPublicKey(pubkeyDecoded)
	if err != nil {
		return err
	}
	pub := pubkey.(*rsa.PublicKey)
	digest := sha3.Sum256(data)
	return rsa.VerifyPSS(pub, crypto.SHA3_256, digest[:], sig, nil)
}
