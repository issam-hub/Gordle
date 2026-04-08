package handlers

import (
	status "gordle-http/internal/handlers/getstatus"
	"gordle-http/internal/handlers/guess"
	"gordle-http/internal/handlers/newgame"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/games", newgame.Handle).Methods(http.MethodPost)
	r.HandleFunc("/games/{id}", status.Handle).Methods(http.MethodGet)
	r.HandleFunc("/games/{id}", guess.Handle).Methods(http.MethodPut)
	return r
}
