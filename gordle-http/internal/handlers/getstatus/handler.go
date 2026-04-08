package status

import (
	"encoding/json"
	"gordle-http/internal/api"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		http.Error(w, "missing the id of the game", http.StatusBadRequest)
		return
	}

	log.Printf("retrieve status of game with id: %v", id)

	apiGame := api.GameResponse{
		ID: id,
	}

	err := json.NewEncoder(w).Encode(apiGame)
	if err != nil {
		log.Printf("failed to write response: %s", err)
	}
}
