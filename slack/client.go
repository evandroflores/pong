package slack

import (
	"context"
	"os"

	"github.com/shomali11/slacker"
	log "github.com/sirupsen/logrus"
)

// Client is the package Client instance
var Client *slacker.Slacker

func init() {
	token := os.Getenv("PONG_TOKEN")
	if token == "" {
		log.Fatal("Bot User OAuth Access Token not found. Set PONG_TOKEN to continue.")
	}

	Client = slacker.NewClient(token)
}

// Listen enables the Client in listen mode
func Listen(ctx context.Context) error {
	return Client.Listen(ctx)
}
