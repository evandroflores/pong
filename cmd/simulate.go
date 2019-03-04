package cmd

import (
	"fmt"

	"github.com/evandroflores/pong/elo"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("simulate <@playerA> vs <@playerB>", "Simulates a game result", simulate)
}

func simulate(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	playerAID := cleanID(request.StringParam("@playerA", ""))
	playerBID := cleanID(request.StringParam("@playerB", ""))

	if !isUser(playerAID) {
		response.ReportError(fmt.Errorf("The given winner is not a User"))
		return
	}

	if !isUser(playerBID) {
		response.ReportError(fmt.Errorf("The given loser is not a User"))
		return
	}

	playerA, errA := model.GetOrCreatePlayer(playerAID)
	if errA != nil {
		response.ReportError(errA)
		return
	}
	playerB, errB := model.GetOrCreatePlayer(playerBID)
	if errB != nil {
		response.ReportError(errB)
		return
	}
	simulateA, simulateB := elo.Calc(playerA.Points, playerB.Points)
	response.Reply(fmt.Sprintf("*%s* wins (%04.f pts) vs *%s* (%04.f pts)", playerA.Name, simulateA, playerB.Name, simulateB))

	simulateB, simulateA = elo.Calc(playerB.Points, playerA.Points)
	response.Reply(fmt.Sprintf("*%s* wins (%04.f pts) vs *%s* (%04.f pts)", playerB.Name, simulateB, playerA.Name, simulateA))
}
