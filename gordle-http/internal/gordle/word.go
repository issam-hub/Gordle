package gordle

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func readFromResponse(response string) string {
	r := regexp.MustCompile(`\["(\w+)"\]`)

	match := r.FindStringSubmatch(string(response))

	return match[1]
}

func GetWord(wordLength int) string {
	fmt.Println("Welcome to Gordle")
	fmt.Println("getting the word for you...")
	response, err := http.Get("https://random-word-api.herokuapp.com/word?length=" + strconv.Itoa(wordLength))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	res, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	word := readFromResponse(string(res))

	fmt.Println("word: ", word)

	return word
}
