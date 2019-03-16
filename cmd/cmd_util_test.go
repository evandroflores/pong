package cmd

import (
	"context"
	"fmt"

	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
	"github.com/shomali11/slacker"
)

type fakeRequest struct {
	ctx        context.Context
	event      *slack.MessageEvent
	properties *proper.Properties
}

func (r *fakeRequest) Param(key string) string {
	return "fake"
}

func (r *fakeRequest) StringParam(key string, defaultValue string) string {
	return "fake"
}

func (r *fakeRequest) BooleanParam(key string, defaultValue bool) bool {
	return false
}

func (r *fakeRequest) IntegerParam(key string, defaultValue int) int {
	return 1
}

func (r *fakeRequest) FloatParam(key string, defaultValue float64) float64 {
	return 1.0
}

func (r *fakeRequest) Context() context.Context {
	return r.ctx
}

func (r *fakeRequest) Event() *slack.MessageEvent {
	return r.event
}

func (r *fakeRequest) Properties() *proper.Properties {
	return r.properties
}

type fakeResponse struct {
	channel  string
	client   *slack.Client
	rtm      *slack.RTM
	messages []string
	errors   []error
}

func (r *fakeResponse) ReportError(err error) {
	r.errors = append(r.errors, err)
}

func (r *fakeResponse) Typing() {
	fmt.Println(r.channel)
}

func (r *fakeResponse) Reply(message string, options ...slacker.ReplyOption) {
	r.messages = append(r.messages, message)
}

func (r *fakeResponse) RTM() *slack.RTM {
	return r.rtm
}

func (r *fakeResponse) Client() *slack.Client {
	return r.client
}

func (r *fakeResponse) GetMessages() []string {
	return r.messages
}

func (r *fakeResponse) GetErrors() []error {
	return r.errors
}
