package cmd

import (
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/assert"

	"github.com/nlopes/slack"
)

var evt = &slack.MessageEvent{
	Msg: slack.Msg{
		Team:    "TTTTTTTTT",
		Channel: "CCCCCCCCC",
		User:    "UUUUUUUUU",
	}}

func TestTryToShowInvalidUser(t *testing.T) {
	var props = proper.NewProperties(
		map[string]string{
			"@someone": "NOTAUSER",
		})

	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	someone(request, response)
	assert.Contains(t, response.GetMessages(), "_Not a User_")
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}

func TestTryToShowTwoUsersNames(t *testing.T) {
	var props = proper.NewProperties(
		map[string]string{
			"@someone": "USER1234USER123",
		})

	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	someone(request, response)
	assert.Contains(t, response.GetMessages(), "_Not a User_")
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}

func TestShowValidUser(t *testing.T) {
	player := model.Player{
		TeamID:    "TTTTTTTTT",
		ChannelID: "CCCCCCCCC",
		SlackID:   "UUUUUUUUU",
		Name:      "Fake User",
	}

	database.Connection.Create(&player)
	defer database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})
	defer database.Connection.Delete(&player)

	var props = proper.NewProperties(
		map[string]string{
			"@someone": "UUUUUUUUU",
		})

	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	someone(request, response)
	assert.Contains(t, response.GetMessages(), "*Fake User* has 1000 points (#01)")
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}
