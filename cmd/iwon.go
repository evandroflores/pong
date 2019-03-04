package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/elo"
	"github.com/evandroflores/udpong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("I won <loser>", "Update your win and the loser points", iWon)
}

func iWon(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := cleanID(request.Event().User)
	loserID := cleanID(request.StringParam("loser", ""))

	if !isUser(loserID) {
		response.ReportError(fmt.Errorf("The given loser is not a User"))
		return
	}

	if winnerID == loserID {
		response.ReportError(fmt.Errorf("Same player? Go find someone to play"))
		return
	}

	winner, _ := model.GetOrCreatePlayer(winnerID)
	loser, _ := model.GetOrCreatePlayer(loserID)

	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	winner.Update()
	loser.Update()

	response.Reply(fmt.Sprintf("*%s* (%04.f pts) x *%s* (%04.f pts)", winner.Name, winner.Points, loser.Name, loser.Points))
}
