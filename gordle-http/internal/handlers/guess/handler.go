package guess

import (
	"encoding/json"
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"log"
	"net/http"
)

type GameGuess interface {
	Modify(id string, r api.GuessRequest) (core.Game, error)
}

func Handler(db GameGuess) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}

		r := api.GuessRequest{}
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		game, err := guess(id, r)

		apiGame := api.ToGameResponse(game)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			log.Printf("failed to write response: %s", err)
		}
	}
}

func guess(id string, r api.GuessRequest) (core.Game, error) {
	return core.Game{
		ID: core.GameID(id),
	}, nil
}
