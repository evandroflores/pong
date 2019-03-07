package cmd

import (
	"fmt"

	"github.com/evandroflores/pong/elo"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("I lost to <@winner>", "Update your lost and the winner points.", iLost)
}

func iLost(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := cleanID(request.StringParam("@winner", ""))
	loserID := cleanID(request.Event().User)

	if !isUser(winnerID) {
		response.ReportError(fmt.Errorf("The given winner is not a User"))
		return
	}

	if winnerID == loserID {
		response.ReportError(fmt.Errorf("Same player? Go find someone to play"))
		return
	}

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	winner, _ := model.GetOrCreatePlayer(teamID, channelID, winnerID)
	loser, _ := model.GetOrCreatePlayer(teamID, channelID, loserID)

	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	winner.Update()
	loser.Update()

	response.Reply(fmt.Sprintf("*%s* (%04.f pts) x *%s* (%04.f pts)", winner.Name, winner.Points, loser.Name, loser.Points))
}
