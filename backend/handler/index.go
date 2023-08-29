package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const endpointIndex = "/index"

func (router *Router) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, r, nil)
	}

	err := timeCheck(r)
	if err != nil {
		fmt.Println(err)
		badRequest(w, r, nil)
		return
	}

	request := Request{
		Method:   "GET",
		Endpoint: endpointIndex,
		User:     r.Header.Get("User"),
		Date:     r.Header.Get("Date"),
		Body:     "",
	}
	_, err = router.verifySignature(&request, r.Header.Get("Signature"))
	if err != nil {
		fmt.Println(err)
		badRequest(w, r, nil)
		return
	}

	index, err := router.db.IndexByEmail(request.User)
	if err != nil {
		fmt.Println(err)
		badRequest(w, r, nil)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(index)
	if err != nil {
		fmt.Println(err)
		internalServerError(w, r, nil)
	}
}
