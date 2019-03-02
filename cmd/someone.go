package cmd

import (
	"fmt"

	"github.com/evandroflores/udpong/model"
	"github.com/shomali11/slacker"
)

func init() {
	Register("show <someone>", "Show someone's points.", someone)
}

func someone(request slacker.Request, response slacker.ResponseWriter) {
	response.Typing()

	userID := cleanID(request.StringParam("someone", ""))

	user, _ := model.GetOrCreatePlayer(userID)
	response.Reply(fmt.Sprintf("%s has %0.f points", user.Name, user.Points))
}
