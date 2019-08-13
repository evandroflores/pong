package cmd

import (
	"bytes"
	"fmt"

	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("top <limit>", "Show the top N players (default 10, limit 20).", top)
}

func top(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)
	limit := request.IntegerParam("limit", 10)

	if limit > 40 {
		response.Reply("Top is limited to 20 players")
		return
	}

	response.Reply(makeTop(request.Event().Channel, model.GetPlayers(teamID, channelID, limit)))
}

func makeTop(uncleanChannelID string, players []model.Player) string {
	var message bytes.Buffer
	if len(players) == 0 {
		return fmt.Sprintf("No rank for channel <#%s>\n\n", uncleanChannelID)
	}
	message.WriteString(fmt.Sprintf("\n*Rank for * <#%s>\n\n", uncleanChannelID))
	for position, player := range players {
		message.WriteString(fmt.Sprintf("*%02d* - %04.f - %s\n", position+1, player.Points, player.Name))
	}

	return message.String()
}
