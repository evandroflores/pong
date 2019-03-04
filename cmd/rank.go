package cmd

import (
	"bytes"
	"fmt"

	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("rank", "Show the top 10 rank.", top10)
	Register("rank all", "Show the entire rank.", all)
}

func top10(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()
	response.Reply(makeRank(model.GetPlayers(10)))
}

func all(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()
	response.Reply(makeRank(model.GetAllPlayers()))
}

func makeRank(players []model.Player) string {
	var message bytes.Buffer
	for position, player := range players {
		message.WriteString(fmt.Sprintf("*%02d* - %04.f - %s\n", position+1, player.Points, player.Name))
	}

	return message.String()
}
