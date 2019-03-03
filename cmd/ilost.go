package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/elo"
	"github.com/evandroflores/udpong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("I lost to <winner>", "Calculate and rank the given Winner and Loser.", iLost)
}

func iLost(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := cleanID(request.StringParam("winner", ""))
	loserID := cleanID(request.Event().User)

	if !isUser(winnerID) {
		response.ReportError(fmt.Errorf("The given winner is not a User"))
		return
	}

	winner, _ := model.GetOrCreatePlayer(winnerID)
	loser, _ := model.GetOrCreatePlayer(loserID)

	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	winner.Update()
	loser.Update()

	response.Reply(fmt.Sprintf("*%s* (%04.f pts) x *%s* (%04.f pts)", winner.Name, winner.Points, loser.Name, loser.Points))
}
