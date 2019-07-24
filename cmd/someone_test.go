package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/nlopes/slack"
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
	fmt.Println(response.GetBlocks())
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
	blocks := response.GetBlocks()

	assert.Len(t, response.GetBlocks(), 1)
	assert.Equal(t, slack.MBTContext, blocks[0].BlockType())
	contextBlock := blocks[0].(*slack.ContextBlock)
	assert.Len(t, contextBlock.ContextElements.Elements, 2)
	potentialImage := contextBlock.ContextElements.Elements[0]
	assert.Equal(t, slack.MixedElementImage, potentialImage.MixedElementType())
	potentialText := contextBlock.ContextElements.Elements[1]
	assert.Equal(t, slack.MixedElementText, potentialText.MixedElementType())
	assert.Equal(t, fmt.Sprintf("*%s* (%04.f pts) *#%02d*", player.Name, 1000.00, 1),
		potentialText.(*slack.TextBlockObject).Text)

	assert.Empty(t, response.GetErrors())
}
