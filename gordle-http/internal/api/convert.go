package api

import "gordle-http/internal/core"

func ToGameResponse(g core.Game) GameResponse {
	var guesses []Guess
	if g.Guesses != nil {
		guesses = make([]Guess, len(g.Guesses))
		for index := 0; index < len(g.Guesses); index++ {
			guesses[index].Word = g.Guesses[index].Word
			guesses[index].Feedback = g.Guesses[index].Feedback
		}
	}

	apiGame := GameResponse{
		ID:           string(g.ID),
		AttemptsLeft: g.AttemptsLeft,
		Guesses:      guesses,
		Status:       string(g.Status),
	}

	if g.AttemptsLeft == 0 {
		apiGame.Solution = ""
	}

	return apiGame
}
