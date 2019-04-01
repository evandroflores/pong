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
	return r.StringParam(key, "")
}

func (r *fakeRequest) StringParam(key, defaultValue string) string {
	return r.properties.StringParam(key, defaultValue)
}

func (r *fakeRequest) BooleanParam(key string, defaultValue bool) bool {
	return r.properties.BooleanParam(key, defaultValue)
}

func (r *fakeRequest) IntegerParam(key string, defaultValue int) int {
	return r.properties.IntegerParam(key, defaultValue)
}

func (r *fakeRequest) FloatParam(key string, defaultValue float64) float64 {
	return r.properties.FloatParam(key, defaultValue)
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
	errors   []string
}

func (r *fakeResponse) ReportError(err error) {
	r.errors = append(r.errors, err.Error())
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

func (r *fakeResponse) GetErrors() []string {
	return r.errors
}
