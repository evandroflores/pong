package cmd

import (
	"fmt"

	"github.com/evandroflores/pong/elo"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("save <@winner> beats <@loser>", "Records points for a given Winner and Loser.", beats)
}

func beats(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := cleanID(request.StringParam("@winner", ""))
	loserID := cleanID(request.StringParam("@loser", ""))

	if !isUser(winnerID) {
		response.Reply("_The given winner is not a User_")
		return
	}

	if !isUser(loserID) {
		response.Reply("_The given loser is not a User_")
		return
	}

	if winnerID == loserID {
		response.Reply("_Same player? Go find someone to play_")
		return
	}

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	winner, errW := model.GetOrCreatePlayer(teamID, channelID, winnerID)
	if errW != nil {
		response.ReportError(errW)
		return
	}
	loser, errL := model.GetOrCreatePlayer(teamID, channelID, loserID)
	if errL != nil {
		response.ReportError(errL)
		return
	}
	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	_ = winner.Update()
	_ = loser.Update()

	response.Reply(fmt.Sprintf("*%s* (%04.f pts) x *%s* (%04.f pts)", winner.Name, winner.Points, loser.Name, loser.Points))
}
