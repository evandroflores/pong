package slack

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
)

const errorFormat = "_%s_"

func NewCustomResponseWriter(channel string, client *slack.Client, rtm *slack.RTM) slacker.ResponseWriter {
	return &CustomResponseWriter{channel: channel, client: client, rtm: rtm}
}

type CustomResponseWriter struct {
	channel string
	client  *slack.Client
	rtm     *slack.RTM
}

func (r *CustomResponseWriter) ReportError(err error) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.channel))
}

func (r *CustomResponseWriter) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.channel))
}

func (r *CustomResponseWriter) Reply(message string, options ...slacker.ReplyOption) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(message, r.channel))
}

func (r *CustomResponseWriter) RTM() *slack.RTM {
	return r.rtm
}

func (r *CustomResponseWriter) Client() *slack.Client {
	return r.client
}
