package cmd

import (
	"github.com/evandroflores/pong/elo"
	"github.com/shomali11/slacker"
)

func init() {
	Register("save <@winner> beats <@loser>", "Records points for a given Winner and Loser.", beats)
}

func beats(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	winnerID := cleanID(request.StringParam("@winner", ""))
	loserID := cleanID(request.StringParam("@loser", ""))

	winner, loser, err := getMatchPlayers(teamID, channelID, winnerID, loserID)
	if err != nil {
		response.ReportError(err)
		return
	}

	winner.Points, loser.Points = elo.Calc(winner.Points, loser.Points)
	winnerDiffPos, _ := winner.Update()
	loserDiffPos, _ := loser.Update()

	response.Reply("", slacker.WithBlocks(versusMessageBlock(&winner, winnerDiffPos, &loser, loserDiffPos)))
}
