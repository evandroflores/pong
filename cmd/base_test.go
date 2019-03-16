package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/model"
	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/assert"
)

const channelIDUnclean string = "<#C12345678>"
const userIDUnclean string = "<@U12345678>"
const teamID string = "T12345678"
const channelID string = "C12345678"
const userID string = "U12345678"

func makeTestEvent() *slack.MessageEvent {
	return &slack.MessageEvent{
		Msg: slack.Msg{
			Team:    "TTTTTTTTT",
			Channel: "CCCCCCCCC",
			User:    "UUUUUUUUU",
		}}
}

func makeTestPlayer() model.Player {
	return model.Player{
		TeamID:    "TTTTTTTTT",
		ChannelID: "CCCCCCCCC",
		SlackID:   "UUUUUUUUU",
		Name:      "Fake User",
	}
}

func TestCleanID(t *testing.T) {
	assert.NotContains(t, cleanID(teamID), "<", ">", "#", "@")
	assert.NotContains(t, cleanID(channelIDUnclean), "<", ">", "#", "@")
	assert.NotContains(t, cleanID(userIDUnclean), "<", ">", "#", "@")
	assert.Equal(t, channelID, cleanID(channelIDUnclean))
	assert.Equal(t, userID, cleanID(userIDUnclean))
}

func TestIsUser(t *testing.T) {
	assert.True(t, isUser(userID))
	assert.False(t, isUser(userIDUnclean))
	assert.False(t, isUser(channelIDUnclean))
	assert.False(t, isUser(channelID))
	assert.False(t, isUser(teamID))
	assert.False(t, isUser(fmt.Sprintf("%s ", userID)))
	assert.False(t, isUser(fmt.Sprintf(" %s ", userID)))
	assert.False(t, isUser(fmt.Sprintf("%s%s", userID, userID)))
	assert.False(t, isUser(fmt.Sprintf("%s %s", userID, userID)))
	assert.False(t, isUser(""))
	assert.False(t, isUser(" "))
}

func TestSayWhat(t *testing.T) {
	var props = proper.NewProperties(map[string]string{})
	evt := makeTestEvent()
	evt.Msg.Text = "Testing"
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	sayWhat(request, response)
	assert.Contains(t, response.GetMessages(), "I have no idea what you mean by _Testing_")
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}
