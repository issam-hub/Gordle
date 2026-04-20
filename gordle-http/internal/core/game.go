package core

import (
	"errors"
)

type Game struct {
	ID           GameID
	Solution     string
	AttemptsLeft byte
	Guesses      []Guess
	Status       Status
}

type GameID string
type Status string

const (
	StatusPlaying = "Playing"
	StatusWon     = "Won"
	StatusLost    = "Lost"
)

type Guess struct {
	Word     string
	Feedback string
}

var ErrGameOver = errors.New("game over")

const MAX_ATTEMPTS = 6
const WORD_LENGTH = 5
