package handler

import (
	"api/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type payloadItem struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Key         string `json:"key"`
	Crendential string `json:"credential"`
}

const endpointItem = "/item"

func (router *Router) Item(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet && r.Method != http.MethodDelete {
		methodNotAllowed(w, r, nil)
		return
	}

	err := timeCheck(r)
	if err != nil {
		fmt.Println(err)
		badRequest(w, r, nil)
		return
	}

	switch r.Method {
	case http.MethodPost:
		router.itemPost(w, r)

	case http.MethodGet:
		router.itemGet(w, r)

	case http.MethodDelete:
		router.itemDelete(w, r)
	}
}

func (router *Router) itemPost(
	w http.ResponseWriter,
	r *http.Request,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		internalServerError(w, r, nil)
		return
	}

	request := Request{
		Method:   "POST",
		Endpoint: endpointItem,
		User:     r.Header.Get("User"),
		Date:     r.Header.Get("Date"),
		Body:     string(body),
	}
	_, err = router.verifySignature(&request, r.Header.Get("Signature"))
	if err != nil {
		badRequest(w, r, nil)
		return
	}

	var payload payloadItem
	err = json.Unmarshal(body, &payload)
	if err != nil {
		badRequest(w, r, nil)
		return
	}

	item := database.Item{
		Email:      request.User,
		Id:         payload.Id,
		Label:      payload.Label,
		Key:        payload.Key,
		Credential: payload.Crendential,
	}
	err = router.db.InsertItem(&item)
	if err != nil {
		internalServerError(w, r, nil)
	}
}

func (router *Router) itemGet(
	w http.ResponseWriter,
	r *http.Request,
) {
	request := Request{
		Method:   "GET",
		Endpoint: endpointItem,
		User:     r.Header.Get("User"),
		Date:     r.Header.Get("Date"),
		Body:     "",
	}
	_, err := router.verifySignature(&request, r.Header.Get("Signature"))
	if err != nil {
		badRequest(w, r, nil)
		return
	}
	id := r.URL.Query().Get("id")
	item, err := router.db.ItemByPK(request.User, id)
	if err != nil {
		badRequest(w, r, nil)
		return
	}
	payload := payloadItem{
		Id:          id,
		Label:       item.Label,
		Key:         item.Key,
		Crendential: item.Credential,
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(payload)
	if err != nil {
		internalServerError(w, r, nil)
	}
}

func (router *Router) itemDelete(
	w http.ResponseWriter,
	r *http.Request,
) {
	request := Request{
		Method:   "DELETE",
		Endpoint: endpointItem,
		User:     r.Header.Get("User"),
		Date:     r.Header.Get("Date"),
		Body:     "",
	}
	_, err := router.verifySignature(&request, r.Header.Get("Signature"))
	if err != nil {
		badRequest(w, r, nil)
		return
	}
	id := r.URL.Query().Get("id")
	err = router.db.DeleteItem(request.User, id)
	if err != nil {
		badRequest(w, r, nil)
		return
	}
}
