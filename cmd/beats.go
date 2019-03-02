package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/slack"
	"github.com/shomali11/slacker"
)

func init() {
	Register("<winner> beats <loser>", "Calculate and rank the given Winner and Loser.", beats)
}

func beats(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := request.StringParam("winner", "")
	loserID := request.StringParam("loser", "")

	winner, _ := slack.Client.GetUserInfo(cleanID(winnerID))
	loser, _ := slack.Client.GetUserInfo(cleanID(loserID))
	response.Reply(fmt.Sprintf("%s x %s", winner.RealName, loser.RealName))
}
