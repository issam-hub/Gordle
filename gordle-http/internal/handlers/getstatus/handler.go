package status

import (
	"encoding/json"
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"log"
	"net/http"
)

type gameGetter interface {
	Get(id string) (core.Game, error)
}

func Handler(db gameGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}

		log.Printf("retrieve status of game with id: %v", id)

		game, err := db.Get(id)
		if err != nil {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}

		apiGame := api.ToGameResponse(game)

		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			log.Printf("failed to write response: %s", err)
		}
	}
}
