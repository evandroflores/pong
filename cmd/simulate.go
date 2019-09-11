package cmd

import (
	"fmt"

	"github.com/evandroflores/pong/elo"
	"github.com/shomali11/slacker"
)

func init() {
	Register("simulate <@playerA> vs <@playerB>", "Simulates a game result.", simulate)
}

func simulate(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	playerAID := cleanID(request.StringParam("@playerA", ""))
	playerBID := cleanID(request.StringParam("@playerB", ""))

	playerA, playerB, err := getMatchPlayers(teamID, channelID, playerAID, playerBID)
	if err != nil {
		response.ReportError(err)
		return
	}

	simulateA, simulateB, _ := elo.Calc(playerA.Points, playerB.Points)
	response.Reply(fmt.Sprintf("*%s* wins (%04.f pts) vs *%s* (%04.f pts)", playerA.Name, simulateA, playerB.Name, simulateB))

	simulateB, simulateA, _ = elo.Calc(playerB.Points, playerA.Points)
	response.Reply(fmt.Sprintf("*%s* wins (%04.f pts) vs *%s* (%04.f pts)", playerB.Name, simulateB, playerA.Name, simulateA))
}
