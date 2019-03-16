package cmd

import (
	"context"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
	"github.com/shomali11/slacker"
	"github.com/stretchr/testify/assert"
)

func TestTryToShowInvalidUser(t *testing.T) {
	defer monkey.UnpatchAll()

	mockMsgOptionUser := func(string) slack.MsgOption {
		t.Fatal("called")
		return nil
	}
	_ = monkey.Patch(slack.MsgOptionUser, mockMsgOptionUser)

	var client = slack.New("FakeClient")
	var rtm = client.NewRTM()
	var mockResponse = slacker.NewResponse("X", client, rtm)

	var info = &slack.Info{
		URL:      "fake",
		User:     &slack.UserDetails{ID: "abc123"},
		Team:     &slack.Team{},
		Users:    []slack.User{},
		Channels: []slack.Channel{},
		Groups:   []slack.Group{},
		Bots:     []slack.Bot{},
		IMs:      []slack.IM{},
	}

	monkey.PatchInstanceMethod(reflect.TypeOf(rtm), "PostMessage",
		func(rtm *slack.RTM, channelID string, options ...slack.MsgOption) (string, string, error) {
			t.Log("Monkey PostMessage")
			assert.Equal(t, "_Not a User_", options[1])
			return "", "", nil
		})

	monkey.PatchInstanceMethod(reflect.TypeOf(rtm), "NewTypingMessage",
		func(rtm *slack.RTM, channel string) *slack.OutgoingMessage {
			return &slack.OutgoingMessage{}
		})

	monkey.PatchInstanceMethod(reflect.TypeOf(rtm), "SendMessage",
		func(rtm *slack.RTM, msg *slack.OutgoingMessage) {
			t.Log("Monkey SendMessage")
		})

	monkey.PatchInstanceMethod(reflect.TypeOf(rtm), "GetInfo",
		func(*slack.RTM) *slack.Info {
			return info
		})

	monkey.PatchInstanceMethod(reflect.TypeOf(client), "ConnectRTM",
		func(*slack.Client) (*slack.Info, string, error) {
			return info, "", nil
		})
	monkey.PatchInstanceMethod(reflect.TypeOf(client), "StartRTM",
		func(*slack.Client) (*slack.Info, string, error) {
			return info, "", nil
		})

	var evt = &slack.MessageEvent{
		Msg: slack.Msg{
			Type:    "",
			Channel: "CCCCCCCCC",
			User:    "UUUUUUUUUU",
			Text:    "TEST",
		}}

	var props = proper.NewProperties(map[string]string{"@someone": "AAA"})
	var mockRequest = slacker.NewRequest(context.Background(), evt, props)

	//client.StartRTM()

	//t.Fatal("*** CHECK ***", slack.MsgOptionUser(rtm.GetInfo().User.ID))

	someone(mockRequest, mockResponse)
}
