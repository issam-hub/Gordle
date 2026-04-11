package newgame

import (
	"encoding/json"
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"log"
	"net/http"
)

type gameAdder interface {
	Add(game core.Game) error
}

func Handler(db gameAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		game, err := createGame()
		if err != nil {
			log.Printf("unable to create a new game: %s", err.Error())
			http.Error(w, "failed to create a new game", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		apiGame := api.ToGameResponse(game)

		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			log.Printf("failed to write response: %s", err)
		}
	}
}

func createGame() (core.Game, error) {
	return core.Game{}, nil
}
