package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

const (
	badRequestMessage          = "Bad Request"
	unauthorizedMessage        = "Unauthorized"
	methodNotAllowedMessage    = "Method Not Allowed"
	basicAuthErrorMessage      = "Malformed Authorization header: Expected Basic Auth"
	internalServerErrorMessage = "Internal Server Error"
)

var (
	errBadrequest          = errors.New(strings.ToLower(badRequestMessage))
	errUnauthorized        = errors.New(strings.ToLower(unauthorizedMessage))
	errMethodNotAllowed    = errors.New(strings.ToLower(methodNotAllowedMessage))
	errBasicAuthError      = errors.New(strings.ToLower(basicAuthErrorMessage))
	errInternalServerError = errors.New(strings.ToLower(internalServerErrorMessage))
)

func logError(r *http.Request, err error) {
	color.Set(color.FgRed)
	defer color.Unset()
	log.Printf("%s: %s %s: %s", r.RemoteAddr, r.Method, r.URL.String(), err)
}

func badRequest(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errBadrequest
	}
	logError(r, err)
	http.Error(w, badRequestMessage, http.StatusBadRequest)
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errUnauthorized
	}
	logError(r, err)
	http.Error(w, unauthorizedMessage, http.StatusUnauthorized)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errMethodNotAllowed
	}
	logError(r, err)
	http.Error(w, methodNotAllowedMessage, http.StatusMethodNotAllowed)
}

func basicAuthError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errBasicAuthError
	}
	logError(r, err)
	http.Error(w, basicAuthErrorMessage, http.StatusBadRequest)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errInternalServerError
	}
	logError(r, err)
	http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
}
