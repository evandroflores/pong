package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/evandroflores/pong/cmd"
	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/slack"
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
		log.Error("Could not start the bot", err)
		runtime.Goexit()
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
