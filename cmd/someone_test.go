package cmd

import (
	"testing"

	"github.com/shomali11/proper"
	"github.com/stretchr/testify/assert"

	"github.com/nlopes/slack"
)

func TestTryToShowInvalidUser(t *testing.T) {
	var evt = &slack.MessageEvent{
		Msg: slack.Msg{
			Team:    "TTTTTTTTT",
			Channel: "CCCCCCCCC",
			User:    "UUUUUUUUU",
			Text:    "show <@UUUUUUUUU>",
		}}

	var props = proper.NewProperties(
		map[string]string{
			"@someone": "NOTAUSER",
		})

	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	someone(request, response)
	assert.Contains(t, response.GetMessages(), "_Not a User_")
}
