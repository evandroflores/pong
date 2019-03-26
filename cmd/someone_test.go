package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/assert"
)

func TestTryToShowInvalidUser(t *testing.T) {
	var props = proper.NewProperties(
		map[string]string{
			"@someone": "NOTAUSER",
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
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

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	someone(request, response)
	assert.Contains(t, response.GetMessages(), "_Not a User_")
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}

func TestTryToShowEmpty(t *testing.T) {
	var props = proper.NewProperties(
		map[string]string{
			"@someone": "",
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	someone(request, response)
	assert.Contains(t, response.GetMessages(), "_Not a User_")
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}

func TestShowValidUser(t *testing.T) {
	player := makeTestPlayer()

	database.Connection.Create(&player)
	defer database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})
	defer database.Connection.Delete(&player)

	var props = proper.NewProperties(
		map[string]string{
			"@someone": player.SlackID,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	someone(request, response)
	assert.Contains(t, response.GetMessages(), fmt.Sprintf("*%s* has 1000 points (#01)", player.Name))
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}
