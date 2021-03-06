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
	elements := contextBlock.ContextElements.Elements
	assert.Len(t, elements, 4)

	assert.Equal(t, slack.MixedElementText, elements[0].MixedElementType())
	assert.Equal(t, fmt.Sprintf("*#%02d*", player.GetPosition()), elements[0].(*slack.TextBlockObject).Text)

	assert.Equal(t, slack.MixedElementImage, elements[1].MixedElementType())
	assert.Equal(t, player.Image, elements[1].(*slack.ImageBlockElement).ImageURL)
	assert.Equal(t, player.Name, elements[1].(*slack.ImageBlockElement).AltText)

	assert.Equal(t, slack.MixedElementText, elements[2].MixedElementType())
	assert.Equal(t, fmt.Sprintf("*%s*", player.Name), elements[2].(*slack.TextBlockObject).Text)

	assert.Equal(t, slack.MixedElementText, elements[3].MixedElementType())
	assert.Equal(t, fmt.Sprintf("(%04.f pts)", player.Points), elements[3].(*slack.TextBlockObject).Text)

	assert.Empty(t, response.GetErrors())
}
