package cmd

import (
	"github.com/shomali11/slacker"
)

func init() {
	Register("I won <@loser>", "Update your win and the loser points.", iWon)
	Register("I beat <@loser>", "Update your win and the loser points.", iWon)
	Register("I crushed <@loser>", "Update your win and the loser points.", iWon)
}

func iWon(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	winnerID := cleanID(request.Event().User)
	loserID := cleanID(request.StringParam("@loser", ""))

	handleMatch(response, teamID, channelID, winnerID, loserID)
}
