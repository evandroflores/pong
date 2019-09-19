package cmd

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/evandroflores/pong/model"
	sl "github.com/evandroflores/pong/slack"
	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/assert"
)

const channelIDUnclean string = "<#C12345678>"
const userIDUnclean string = "<@U12345678>"
const teamID string = "T12345678"
const channelID string = "C12345678"
const userID string = "U12345678"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func makeTestEvent() *slack.MessageEvent {
	return &slack.MessageEvent{
		Msg: slack.Msg{
			Team:    "TTTTTTTTT",
			Channel: "CCCCCCCCC",
			User:    "UUUUUUUUU",
		}}
}

func makeTestPlayer() model.Player {
	randomInt := rand.Intn(10000000)
	name := fmt.Sprintf("Fake User - %08d", randomInt)
	slackID := fmt.Sprintf("U%08d", randomInt)

	mockIngest := func(player *model.Player) {
		fmt.Printf("Ingesting called for %s\n", player.ToStr())
	}
	player := model.Player{
		TeamID:    "TTTTTTTTT",
		ChannelID: "CCCCCCCCC",
		SlackID:   slackID,
		Name:      name,
		Image:     "https://www.thispersondoesnotexist.com/image",
	}
	monkey.PatchInstanceMethod(reflect.TypeOf(&player), "IngestData", mockIngest)

	return player
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

func TestInvalidWinner(t *testing.T) {
	var props = proper.NewProperties(
		map[string]string{
			"@winner": "INVALID",
		})
	for _, command := range sl.Client.BotCommands() {
		if !strings.Contains(command.Usage(), "@winner") {
			continue
		}
		request := &fakeRequest{event: makeTestEvent(), properties: props}
		response := &fakeResponse{}

		command.Execute(request, response)

		assert.Contains(t, response.GetErrors(), "the given winner is not a user")
		assert.Len(t, response.GetErrors(), 1)
	}
}

func TestInvalidLoser(t *testing.T) {
	var props = proper.NewProperties(
		map[string]string{
			"@winner": "U00000000", // Needs a 'valid' winner for cmd.beats
			"@loser":  "INVALID",
		})
	for _, command := range sl.Client.BotCommands() {
		if !strings.Contains(command.Usage(), "@loser") {
			continue
		}
		request := &fakeRequest{event: makeTestEvent(), properties: props}
		response := &fakeResponse{}

		command.Execute(request, response)

		assert.Contains(t, response.GetErrors(), "the given loser is not a user")
		assert.Len(t, response.GetErrors(), 1)
	}
}

func TestForeverAlone(t *testing.T) {
	evt := makeTestEvent()
	var props = proper.NewProperties(
		map[string]string{
			"@winner": evt.User,
			"@loser":  evt.User,
		})

	for _, command := range sl.Client.BotCommands() {
		if !strings.HasPrefix(command.Usage(), "I ") &&
			!strings.Contains(command.Usage(), " beats ") {
			continue
		}
		request := &fakeRequest{event: makeTestEvent(), properties: props}
		response := &fakeResponse{}

		command.Execute(request, response)

		assert.Contains(t, response.GetErrors(), "go find someone to play")
		assert.Len(t, response.GetErrors(), 1)
	}
}

func TestTrend(t *testing.T) {
	assert.Equal(t, getPosDiff(0), "")
	assert.Equal(t, getPosDiff(1), " ↑ 1 ")
	assert.Equal(t, getPosDiff(-1), " ↓ 1 ")
	assert.Equal(t, getPosDiff(5), " ↑ 5 ")
	assert.Equal(t, getPosDiff(-5), " ↓ 5 ")
}
