package cmd

import (
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

	if limit > 20 {
		response.Reply("Top is limited to 20 players")
		return
	}

	players := model.GetPlayers(teamID, channelID, limit)
	if len(players) == 0 {
		response.Reply(fmt.Sprintf("No rank for channel <#%s>", channelID))
		return
	}
	header := fmt.Sprintf("*Rank for * <#%s>", channelID)

	response.Reply("", slacker.WithBlocks(listMessageBlock(header, players)))

}
