package newgame

import (
	"encoding/json"
	"fmt"
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"gordle-http/internal/gordle"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type gameAdder interface {
	Add(game core.Game) error
}

func Handler(db gameAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		game, err := createGame(req.Body, db)
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

func createGame(reader io.Reader, db gameAdder) (core.Game, error) {
	solution := gordle.GetWord(core.WORD_LENGTH)
	game := gordle.New(reader, solution, core.MAX_ATTEMPTS)

	g := core.Game{
		ID:           core.GameID(uuid.New().String()),
		Gordle:       *game,
		AttemptsLeft: core.MAX_ATTEMPTS,
		Guesses:      []core.Guess{},
		Status:       core.StatusPlaying,
	}

	err := db.Add(g)
	if err != nil {
		return core.Game{}, fmt.Errorf("failed to save the new game")
	}
	return g, nil
}
