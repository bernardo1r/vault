package handler

import (
	"api/crypto"
	"fmt"
	"net/http"
)

func (router *Router) Login(w http.ResponseWriter, r *http.Request) {
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

	user, err := router.db.UserByEmail(username)
	if err != nil {
		badRequest(w, r, err)
		return
	}

	_, err = crypto.GCMDecrypt(passwordDecoded, nil, user.Totp)
	if err != nil {
		unauthorized(w, r, fmt.Errorf("decrypting totp secret: %w", err))
		return
	}
}
