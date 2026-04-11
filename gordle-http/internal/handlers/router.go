package handlers

import (
	status "gordle-http/internal/handlers/getstatus"
	"gordle-http/internal/handlers/guess"
	"gordle-http/internal/handlers/newgame"
	"gordle-http/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
)

func Router(db *repository.GameRepository) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/games", newgame.Handler(db)).Methods(http.MethodPost)
	r.HandleFunc("/games/{id}", status.Handler(db)).Methods(http.MethodGet)
	r.HandleFunc("/games/{id}", guess.Handler(db)).Methods(http.MethodPut)
	return r
}
