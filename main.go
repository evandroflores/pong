package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	token := os.Getenv("PONG_TOKEN")
	if token == "" {
		log.Fatal("Bot User OAuth Access Token not found. Set PONG_TOKEN to continue.")
	}
}
