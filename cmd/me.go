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

	user, _ := model.GetOrCreatePlayer(userID)
	response.Reply(fmt.Sprintf("You have %0.f points", user.Points))
}
