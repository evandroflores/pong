package cmd

import (
	"fmt"

	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("show <@someone>", "Show someone's points.", someone)
}

func someone(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	userID := cleanID(request.StringParam("@someone", ""))
	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	if !isUser(userID) {
		response.Reply("_Not a User_")
		return
	}

	user, err := model.GetOrCreatePlayer(teamID, channelID, userID)
	if err != nil {
		response.ReportError(err)
		return
	}

	blocks := user.GetBlockCard(fmt.Sprintf("(%04.f pts) #%d", user.Points, user.GetPosition()))
	response.Reply("", slacker.WithBlocks(blocks))
}
