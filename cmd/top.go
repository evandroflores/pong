package cmd

import (
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("top <limit>", "Show the top N players (default 10).", top)
}

func top(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)
	limit := request.IntegerParam("limit", 10)

	response.Reply(makeRank(request.Event().Channel, model.GetPlayers(teamID, channelID, limit)))
}
