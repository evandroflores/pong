package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/elo"
	"github.com/evandroflores/udpong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("<@winner> beats <@loser>", "Records points for a given Winner and Loser.", beats)
}

func beats(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := cleanID(request.StringParam("@winner", ""))
	loserID := cleanID(request.StringParam("@loser", ""))

	if !isUser(winnerID) {
		response.ReportError(fmt.Errorf("The given winner is not a User"))
		return
	}

	if !isUser(loserID) {
		response.ReportError(fmt.Errorf("The given loser is not a User"))
		return
	}

	if winnerID == loserID {
		response.ReportError(fmt.Errorf("Same player? Go find someone to play"))
		return
	}

	winner, errW := model.GetOrCreatePlayer(winnerID)
	if errW != nil {
		response.ReportError(errW)
		return
	}
	loser, errL := model.GetOrCreatePlayer(loserID)
	if errL != nil {
		response.ReportError(errL)
		return
	}
	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	winner.Update()
	loser.Update()

	response.Reply(fmt.Sprintf("*%s* (%04.f pts) x *%s* (%04.f pts)", winner.Name, winner.Points, loser.Name, loser.Points))
}
