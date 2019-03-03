package cmd

import (
	"fmt"
	"strings"

	"github.com/evandroflores/udpong/slack"
	"github.com/shomali11/slacker"
	log "github.com/sirupsen/logrus"
)

var replacer = strings.NewReplacer("<", "", ">", "", "@", "")

const userPrefix = "U"

// Register add a command to commands list an prepare to register to slacker
func Register(usage string, description string, handler func(request slacker.Request, response slacker.ResponseWriter)) {
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

func isUser(SlackID string) bool {
	if strings.HasPrefix(SlackID, userPrefix) {
		return true
	}
	return false
}
