package handler

import (
	"api/database"
	"log"
	"net/http"
)

type Router struct {
	db database.DB
}

func New(db database.DB) *Router {
	return &Router{
		db: db,
	}
}

func Log(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logIp(r)
		handler(w, r)
	}
}

func logIp(r *http.Request) {
	log.Printf("%s: %s %s", r.RemoteAddr, r.Method, r.URL.String())
}
