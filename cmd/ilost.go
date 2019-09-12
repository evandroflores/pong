package cmd

import (
	"github.com/shomali11/slacker"
)

func init() {
	Register("I lost to <@winner>", "Update your lost and the winner points.", iLost)
}

func iLost(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	winnerID := cleanID(request.StringParam("@winner", ""))
	loserID := cleanID(request.Event().User)

	handleMatch(response, teamID, channelID, winnerID, loserID)
}
