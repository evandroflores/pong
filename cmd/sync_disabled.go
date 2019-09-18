package cmd

import (
	"bytes"
	"fmt"

	"github.com/evandroflores/pong/model"
	"github.com/evandroflores/pong/slack"
	"github.com/shomali11/slacker"
)

func init() {
	RegisterAdmin("sync-disabled", "Sync disabled Slack users and purge from Pong", syncDisabled)
}

func syncDisabled(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	players := model.GetAllPlayers(teamID, channelID)

	var message bytes.Buffer

	removed := 0
	for _, player := range players {
		slackUser, _ := slack.Client.GetUserInfo(player.SlackID)
		if slackUser.Deleted {
			player.Delete()
			message.WriteString(fmt.Sprintf("*%s* - Removed\n", player.Name))
			removed++
		}
	}
	if removed == 0 {
		message.WriteString("No users removed.")
	}

	response.Reply(message.String())
}
