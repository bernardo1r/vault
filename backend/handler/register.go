package handler

import (
	"api/crypto"
	"api/database"
	"encoding/json"
	"fmt"
	"net/http"
)

type payloadRegister struct {
	PublicKeyEnc   string `json:"publicKeyEnc"`
	PrivateKeyEnc  string `json:"privateKeyEnc"`
	PublicKeySign  string `json:"publicKeySign"`
	PrivateKeySign string `json:"privateKeySign"`
	Totp           string `json:"totp"`
}

func (router *Router) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	decoder := json.NewDecoder(r.Body)
	var payload payloadRegister
	err = decoder.Decode(&payload)
	if err != nil {
		badRequest(w, r, fmt.Errorf("decoding payload: %w", err))
		return
	}

	var user database.User
	err = fillUserPayload(&user, &payload)
	if err != nil {
		badRequest(w, r, err)
		return
	}

	totp, err := crypto.GCMEncrypt(passwordDecoded, nil, []byte(payload.Totp))
	if err != nil {
		badRequest(w, r, fmt.Errorf("encrypting totp secret: %w", err))
		return
	}

	fillUserCred(&user, username, totp)
	err = router.db.InsertUser(&user)
	if err != nil {
		if database.IsUniqueViolation(err) {
			err = fmt.Errorf("user %s already in database: %w", username, err)
		}
		badRequest(w, r, fmt.Errorf("user insertion in database: %w", err))
	}
}

func fillUserPayload(user *database.User, payload *payloadRegister) error {
	var err error
	user.PublicKeyEnc = payload.PublicKeyEnc
	err = checkEncodedValue(user.PublicKeyEnc)
	if err != nil {
		return fmt.Errorf("public key encryption: %w", err)
	}

	user.PrivateKeyEnc = payload.PrivateKeyEnc
	err = checkEncodedValue(user.PrivateKeyEnc)
	if err != nil {
		return fmt.Errorf("private key encryption: %w", err)
	}

	user.PublicKeySign = payload.PublicKeySign
	err = checkEncodedValue(user.PublicKeySign)
	if err != nil {
		return fmt.Errorf("public key sign: %w", err)
	}
	err = crypto.CheckPublicKey(user.PublicKeySign)
	if err != nil {
		return fmt.Errorf("public key sign: %w", err)
	}

	user.PrivateKeySign = payload.PrivateKeySign
	err = checkEncodedValue(user.PrivateKeySign)
	if err != nil {
		return fmt.Errorf("private key sign: %w", err)
	}
	return nil
}

func fillUserCred(
	user *database.User,
	email string,
	totp []byte,
) {
	user.Email = email
	user.Totp = totp
}
