package status

import (
	"encoding/json"
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"log"
	"net/http"
)

type gameGetter interface {
	Get(id string) core.Game
}

func Handler(db gameGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}

		log.Printf("retrieve status of game with id: %v", id)

		game := getGame(id)

		apiGame := api.ToGameResponse(game)

		err := json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			log.Printf("failed to write response: %s", err)
		}
	}

}
func getGame(id string) core.Game {
	return core.Game{
		ID: core.GameID(id),
	}
}
