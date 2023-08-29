package handler

import (
	"api/crypto"
	"api/database"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Request struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Date     string `json:"date"`
	Body     string `json:"body"`
}

func (router *Router) verifySignature(request *Request, sig string) (*database.User, error) {
	user, err := router.db.UserByEmail(request.User)
	if err != nil {
		return nil, errors.New("user not found")
	}
	sigDecoded, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	err = crypto.VerifySignature(user.PublicKeySign, data, sigDecoded)

	return &user, err
}
