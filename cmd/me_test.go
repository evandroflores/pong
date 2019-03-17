package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"

	"github.com/shomali11/proper"
	"github.com/stretchr/testify/assert"
)

func TestMeNewUser(t *testing.T) {
	evt := makeTestEvent()
	evt.Msg.User = "U00000NEW"
	evt.Msg.Channel = "C00000NEW"

	evtPlayer := &model.Player{
		TeamID:    evt.Msg.Team,
		ChannelID: evt.Msg.Channel,
		SlackID:   evt.Msg.User,
	}
	userFromDb := model.Player{}
	database.Connection.Where(evtPlayer).First(&userFromDb)

	assert.Equal(t, model.Player{}, userFromDb)

	var props = proper.NewProperties(map[string]string{})

	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	me(request, response)
	database.Connection.Where(evtPlayer).Delete(&model.Player{})
	database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})

	assert.Contains(t, response.GetMessages(), fmt.Sprintf("You have 1000 points (#01) on <#%s>", evt.Channel))
	assert.Len(t, response.GetMessages(), 1)
	assert.Empty(t, response.GetErrors())
}
