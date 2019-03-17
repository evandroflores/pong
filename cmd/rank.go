package cmd

import (
	"bytes"
	"fmt"

	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("rank", "Show the entire rank.", rank)
	Register("top <limit>", "Show the top N players (default 10).", top)
}

func rank(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	response.Reply(makeRank(request.Event().Channel, model.GetAllPlayers(teamID, channelID)))
}

func top(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)
	limit := request.IntegerParam("limit", 10)

	response.Reply(makeRank(request.Event().Channel, model.GetPlayers(teamID, channelID, limit)))
}

func makeRank(uncleanChannelID string, players []model.Player) string {
	var message bytes.Buffer
	if len(players) == 0 {
		return "No rank for this channel"
	}
	message.WriteString(fmt.Sprintf("\n*Rank for * <#%s>\n\n", uncleanChannelID))
	for position, player := range players {
		message.WriteString(fmt.Sprintf("*%02d* - %04.f - %s\n", position+1, player.Points, player.Name))
	}

	return message.String()
}
