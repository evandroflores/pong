package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/evandroflores/udpong/cmd"
	"github.com/evandroflores/udpong/database"
	"github.com/evandroflores/udpong/slack"
)

func main() {
	log.SetLevel(log.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())

	killListener()
	cmd.LoadCommands()
	defer cancel()
	defer database.Close()
	err := slack.Listen(ctx)
	if err != nil {
		log.Fatal("Could not start the bot", err)
	}
}

func killListener() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\r")
		database.Close()
		fmt.Printf("\rBye bye.\n\n👋\n\n")
		os.Exit(0)
	}()
}
