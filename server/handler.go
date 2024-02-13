package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	mux *mux.Router
}

func NewHandlers() *Handler {
	h := &Handler{
		mux: mux.NewRouter(),
	}
	h.mux.HandleFunc("/health", h.health)
	return h
}

func (h Handler) Router() *mux.Router {
	return h.mux
}

func (h Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
