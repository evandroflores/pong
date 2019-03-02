package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("<winner> beats <loser>", "Calculate and rank the given Winner and Loser.", beats)
}

func beats(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := cleanID(request.StringParam("winner", ""))
	loserID := cleanID(request.StringParam("loser", ""))

	winner, _ := model.GetOrCreatePlayer(winnerID)
	loser, _ := model.GetOrCreatePlayer(loserID)
	response.Reply(fmt.Sprintf("%s x %s", winner.Name, loser.Name))
}
