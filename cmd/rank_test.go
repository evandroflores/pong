package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/assert"
)

func TestEmptyRankForChannel(t *testing.T) {
	var props = proper.NewProperties(map[string]string{})

	evt := makeTestEvent()
	evt.Msg.Channel = "CAAAAAAAA"
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	rank(request, response)
	assert.Contains(t, response.GetMessages(), "No rank for this channel")
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}

func TestRankForChannel(t *testing.T) {
	defer database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})

	expected := "\n*Rank for * <#CCCCCCCCC>\n\n"
	for i := 1; i <= 20; i++ {
		player := makeTestPlayer()
		player.Name = fmt.Sprintf("%s - %02d", player.Name, i)
		player.SlackID = fmt.Sprintf("UUUUUUU%02d", i)
		player.Points = 1000 - float64(i)
		expected += fmt.Sprintf("*%02d* - %04.f - %s\n", i, player.Points, player.Name)
		database.Connection.Create(&player)
		defer database.Connection.Delete(&player)
	}

	var props = proper.NewProperties(map[string]string{})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	rank(request, response)
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
	assert.Equal(t, expected, response.GetMessages()[0])
}
