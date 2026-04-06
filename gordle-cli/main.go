package main

import (
	"gordle-cli/gordle"
	"os"
)

const MAX_ATTEMPTS = 6
const WORD_LENGTH = 5

func main() {
	solution := gordle.GetWord(WORD_LENGTH)
	g := gordle.New(os.Stdin, solution, MAX_ATTEMPTS)

	g.Play()

}
