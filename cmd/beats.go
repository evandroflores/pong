package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/elo"
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

	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	winner.Update()
	loser.Update()

	response.Reply(fmt.Sprintf("*%s* (%04.f pts) x *%s* (%04.f pts)", winner.Name, winner.Points, loser.Name, loser.Points))
}
