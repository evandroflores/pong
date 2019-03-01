package main

import (
	"context"
	"os"

	"github.com/shomali11/slacker"
	log "github.com/sirupsen/logrus"
)

func main() {
	token := os.Getenv("PONG_TOKEN")
	if token == "" {
		log.Fatal("Bot User OAuth Access Token not found. Set PONG_TOKEN to continue.")
	}

	bot := slacker.NewClient(token)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal("Could not start the bot", err)
	}
}
