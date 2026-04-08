package main

import (
	"gordle-http/internal/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	srv := &http.Server{
		Handler:      handlers.Router(),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("server running at :8000")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
