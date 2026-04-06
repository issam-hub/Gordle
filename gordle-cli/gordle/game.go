package gordle

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

var errInvalidWordLength = fmt.Errorf("invalid guess, word doesn't have the same number of characters as the solution")

func splitToUpperCaseCharacters(word string) []rune {
	return []rune(strings.ToUpper(word))
}

type Game struct {
	reader      *bufio.Reader
	solution    []rune
	maxAttempts int
}

func New(playerInput io.Reader, solution string, maxAttempts int) *Game {
	return &Game{
		reader:      bufio.NewReader(playerInput),
		solution:    splitToUpperCaseCharacters(solution),
		maxAttempts: maxAttempts,
	}
}

func (g *Game) Play() {
	for currentAttempt := 1; currentAttempt <= g.maxAttempts; currentAttempt++ {
		guess := g.ask()

		if slices.Equal(guess, g.solution) {
			fmt.Printf("🎉 You won! You found it in %d guess(es)! The word was: %s\n", currentAttempt, string(g.solution))
			return
		} else {
			out := computeFeedback(guess, g.solution)
			fmt.Println(out.StringConcat())
		}

	}
	fmt.Printf("😞 You've lost! The solution was: %s \n", string(g.solution))
}

func (g *Game) ask() []rune {
	fmt.Printf("Enter a %d-character guess: ", len(g.solution))
	for {
		playerInput, _, err := g.reader.ReadLine()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Gordle failed to read your guess: %s\n", err.Error())
			continue
		}
		guess := splitToUpperCaseCharacters(string(playerInput))

		if err := g.validateGuess(guess); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Your attempt is invalid with Gordle's solution! %s\n", err.Error())
		} else {
			return guess
		}
	}
}

func (g *Game) validateGuess(guess []rune) error {
	if len(guess) != len(g.solution) {
		return fmt.Errorf("expected %d, got %d, %w", len(g.solution), len(guess), errInvalidWordLength)
	}
	return nil
}

func computeFeedback(guess, solution []rune) feedback {
	result := make(feedback, len(guess))
	seenChars := make([]bool, len(solution))

	if len(guess) != len(solution) {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error! Guess and solution"+
			" have different lengths: %d vs %d", len(guess), len(solution))
		return result
	}

	for pos, char := range guess {
		if char == solution[pos] {
			result[pos] = CORRECT_POSITION
			seenChars[pos] = true
		}
	}

	for posG, char := range guess {
		if result[posG] != ABSENT_CHARACTER {
			continue
		}

		for posS, target := range solution {
			if seenChars[posS] {
				continue
			}

			if char == target {
				result[posG] = WRONG_POSITION
				seenChars[posS] = true
				break
			}
		}
	}

	return result
}
