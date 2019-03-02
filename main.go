package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/evandroflores/udpong/cmd"
	"github.com/evandroflores/udpong/database"
	"github.com/evandroflores/udpong/slack"
)

func main() {
	log.SetLevel(log.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())

	cmd.LoadCommands()
	defer cancel()
	defer database.Close()
	err := slack.Listen(ctx)
	if err != nil {
		log.Fatal("Could not start the bot", err)
	}
}
