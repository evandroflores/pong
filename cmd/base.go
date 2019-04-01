package cmd

import (
	"fmt"
	"strings"

	"github.com/evandroflores/pong/model"
	"github.com/evandroflores/pong/slack"
	"github.com/shomali11/slacker"
	log "github.com/sirupsen/logrus"
)

var replacer = strings.NewReplacer("<", "", ">", "", "@", "", "#", "")

const userPrefix = "U"

// Register add a command to commands list an prepare to register to slacker
func Register(usage, description string, handler func(request slacker.Request, response slacker.ResponseWriter)) {
	log.Infof("Registering %s - %s", usage, description)
	slack.Client.Command(usage, &slacker.CommandDefinition{Description: description, Handler: handler})
}

func cleanID(userID string) string {
	return replacer.Replace(userID)
}

// LoadCommands will force `init` all classes on this package
func LoadCommands() {
	slack.Client.DefaultCommand(sayWhat)
	log.Infof("%d commands loaded", len(slack.Client.BotCommands()))
}

func sayWhat(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()
	response.Reply(fmt.Sprintf("I have no idea what you mean by _%s_", request.Event().Text))
}

func isUser(slackID string) bool {
	if len(slackID) > 9 {
		return false
	}
	if strings.HasPrefix(slackID, userPrefix) {
		return true
	}
	return false
}

func getMatchPlayers(teamID, channelID, winnerID, loserID string) (winner, loser model.Player, err error) {
	if !isUser(winnerID) {
		return model.Player{}, model.Player{}, fmt.Errorf("the given winner is not a user")
	}

	if !isUser(loserID) {
		return model.Player{}, model.Player{}, fmt.Errorf("the given loser is not a user")
	}

	if winnerID == loserID {
		return model.Player{}, model.Player{}, fmt.Errorf("go find someone to play")
	}

	winner, errW := model.GetOrCreatePlayer(teamID, channelID, winnerID)
	if errW != nil {
		return model.Player{}, model.Player{}, errW
	}

	loser, errL := model.GetOrCreatePlayer(teamID, channelID, loserID)
	if errL != nil {
		return model.Player{}, model.Player{}, errL
	}

	return winner, loser, nil
}
