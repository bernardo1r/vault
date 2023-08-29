package handler

import (
	"api/crypto"
	"api/database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pquerna/otp/totp"
)

type payloadAuthenticator struct {
	PublicKeyEnc   string `json:"publicKeyEnc"`
	PrivateKeyEnc  string `json:"privateKeyEnc"`
	PublicKeySign  string `json:"publicKeySign"`
	PrivateKeySign string `json:"privateKeySign"`
}

func (router *Router) Authenticator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, r, nil)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		basicAuthError(w, r, nil)
		return
	}
	passwordDecoded, err := crypto.DecodePassword(password)
	if err != nil {
		badRequest(w, r, fmt.Errorf("decoding password: %w", err))
		return
	}

	passcode := r.Header.Get("Passcode")
	if len(passcode) == 0 {
		badRequest(w, r, errors.New("missing 'Passcode' header"))
		return
	}

	user, err := router.db.UserByEmail(username)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	secret, err := crypto.GCMDecrypt(passwordDecoded, nil, user.Totp)
	if err != nil {
		unauthorized(w, r, fmt.Errorf("decrypting totp secret: %w", err))
		return
	}

	if !totp.Validate(passcode, string(secret)) {
		unauthorized(w, r, errors.New("validating totp code"))
		return
	}

	var payload payloadAuthenticator
	fillAuthenticatorPayload(&payload, &user)
	payloadEncoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		internalServerError(w, r, nil)
		return
	}

	_, err = w.Write(payloadEncoded)
	if err != nil {
		internalServerError(w, r, nil)
	}
}

func fillAuthenticatorPayload(payload *payloadAuthenticator, user *database.User) {
	payload.PublicKeyEnc = user.PublicKeyEnc
	payload.PrivateKeyEnc = user.PrivateKeyEnc
	payload.PublicKeySign = user.PublicKeySign
	payload.PrivateKeySign = user.PrivateKeySign
}
