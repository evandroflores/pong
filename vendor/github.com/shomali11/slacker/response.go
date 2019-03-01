package slacker

import (
	"fmt"

	"github.com/nlopes/slack"
)

const (
	errorFormat = "*Error:* _%s_"
)

// A ResponseWriter interface is used to respond to an event
type ResponseWriter interface {
	Reply(text string, options ...ReplyOption)
	ReportError(err error)
	Typing()
	RTM() *slack.RTM
	Client() *slack.Client
}

// NewResponse creates a new response structure
func NewResponse(channel string, client *slack.Client, rtm *slack.RTM) ResponseWriter {
	return &response{channel: channel, client: client, rtm: rtm}
}

type response struct {
	channel string
	client  *slack.Client
	rtm     *slack.RTM
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *response) ReportError(err error) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.channel))
}

// Typing send a typing indicator
func (r *response) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.channel))
}

// Reply send a attachments to the current channel with a message
func (r *response) Reply(message string, options ...ReplyOption) {
	defaults := newReplyDefaults(options...)

	r.rtm.PostMessage(
		r.channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionUser(r.rtm.GetInfo().User.ID),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionAttachments(defaults.Attachments...),
	)
}

// RTM returns the RTM client
func (r *response) RTM() *slack.RTM {
	return r.rtm
}

// Client returns the slack client
func (r *response) Client() *slack.Client {
	return r.client
}
