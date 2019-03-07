package cmd

import (
	"fmt"

	"github.com/evandroflores/pong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("me", "Show your points.", me)
}

func me(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	userID := cleanID(request.Event().User)
	teamID := cleanID(request.Event().Team)
	channelID := cleanID(request.Event().Channel)

	user, err := model.GetOrCreatePlayer(teamID, channelID, userID)
	if err != nil {
		response.ReportError(err)
		return
	}

	response.Reply(fmt.Sprintf("You have %0.f points", user.Points))
}
