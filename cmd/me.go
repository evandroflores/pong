package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("me", "Show your points.", me)
}

func me(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	userID := cleanID(request.Event().User)

	user, err := model.GetOrCreatePlayer(userID)
	if err != nil {
		response.ReportError(err)
		return
	}

	response.Reply(fmt.Sprintf("You have %0.f points", user.Points))
}
