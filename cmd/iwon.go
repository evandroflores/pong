package cmd

import (
	"fmt"

	"github.com/evandroflores/pong/elo"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("I won <@loser>", "Update your win and the loser points.", iWon)
}

func iWon(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	winnerID := cleanID(request.Event().User)
	loserID := cleanID(request.StringParam("@loser", ""))

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

	winner, _ := model.GetOrCreatePlayer(teamID, channelID, winnerID)
	loser, _ := model.GetOrCreatePlayer(teamID, channelID, loserID)

	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	_ = winner.Update()
	_ = loser.Update()

	response.Reply(fmt.Sprintf("*%s* (%04.f pts) x *%s* (%04.f pts)", winner.Name, winner.Points, loser.Name, loser.Points))
}
