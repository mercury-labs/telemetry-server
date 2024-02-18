package server

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
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
	h.mux.HandleFunc("/track", h.track)
	return h
}

func (h Handler) Router() *mux.Router {
	return h.mux
}

func (h Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h Handler) track(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	spew.Dump(m)
	w.WriteHeader(http.StatusOK)
}
