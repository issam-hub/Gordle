package handlers

import (
	status "gordle-http/internal/handlers/getstatus"
	"gordle-http/internal/handlers/guess"
	"gordle-http/internal/handlers/newgame"
	"gordle-http/internal/repository"
	"net/http"
)

func Router(db *repository.GameRepository) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("POST /games", newgame.Handler(db))
	r.HandleFunc("GET /games/{id}", status.Handler(db))
	r.HandleFunc("PUT /games/{id}", guess.Handler(db))
	return r
}
