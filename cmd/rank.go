package cmd

import (
	"bytes"
	"fmt"

	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("rank", "Show the entire rank.", rank)
}

func rank(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	response.Reply(makeRank(request.Event().Channel, model.GetAllPlayers(teamID, channelID)))
}

func makeRank(uncleanChannelID string, players []model.Player) string {
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
