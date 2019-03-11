package cmd

import (
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("top", "Show the top 10.", top10)
}

func top10(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	response.Reply(makeRank(model.GetPlayers(teamID, channelID, 10)))
}
