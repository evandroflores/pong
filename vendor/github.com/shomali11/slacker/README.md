# slacker [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/slacker)](https://goreportcard.com/report/github.com/shomali11/slacker) [![GoDoc](https://godoc.org/github.com/shomali11/slacker?status.svg)](https://godoc.org/github.com/shomali11/slacker) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Built on top of the Slack API [github.com/nlopes/slack](https://github.com/nlopes/slack) with the idea to simplify the Real-Time Messaging feature to easily create Slack Bots, assign commands to them and extract parameters.

## Features

* Easy definitions of commands and their input
* Available bot initialization, errors and default handlers
* Simple parsing of String, Integer, Float and Boolean parameters
* Contains support for `context.Context`
* Built-in `help` command
* Supports authorization
* Bot responds to mentions and direct messages
* Handlers run concurrently via goroutines
* Full access to the Slack API [github.com/nlopes/slack](https://github.com/nlopes/slack)

## Dependencies

* `commander` [github.com/shomali11/commander](https://github.com/shomali11/commander)
* `slack` [github.com/nlopes/slack](https://github.com/nlopes/slack)

# Examples

## Example 1

Defining a command using slacker

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	}

	bot.Command("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 2

Defining a command with an optional description and example

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Ping!",
		Example:     "ping",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	}

	bot.Command("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 3

Defining a command with a parameter

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Echo a word!",
		Example:     "echo hello",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")
			response.Reply(word)
		},
	}

	bot.Command("echo <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 4

Defining a command with two parameters. Parsing one as a string and the other as an integer. 
_(The second parameter is the default value in case no parameter was passed or could not parse the value)_

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Repeat a word a number of times!",
		Example:     "repeat hello 10",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			word := request.StringParam("word", "Hello!")
			number := request.IntegerParam("number", 1)
			for i := 0; i < number; i++ {
				response.Reply(word)
			}
		},
	}

	bot.Command("repeat <word> <number>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 5

Send an error message to the Slack channel

```go
package main

import (
	"context"
	"errors"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Tests errors",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.ReportError(errors.New("Oops!"))
		},
	}

	bot.Command("test", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 6

Send a "Typing" indicator

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
	"time"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Server time!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Typing()

			time.Sleep(time.Second)

			response.Reply(time.Now().Format(time.RFC1123))
		},
	}

	bot.Command("time", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 7

Showcasing the ability to access the [github.com/nlopes/slack](https://github.com/nlopes/slack) API and the Real-Time Messaging Protocol.
_In this example, we are sending a message using RTM and uploading a file using the Slack API._

```go
package main

import (
	"context"
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Upload a word!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")
			channel := request.Event().Channel

			rtm := response.RTM()
			client := response.Client()

			rtm.SendMessage(rtm.NewOutgoingMessage("Uploading file ...", channel))
			client.UploadFile(slack.FileUploadParameters{Content: word, Channels: []string{channel}})
		},
	}

	bot.Command("upload <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 8

Showcasing the ability to leverage `context.Context` to add a timeout

```go
package main

import (
	"context"
	"errors"
	"github.com/shomali11/slacker"
	"log"
	"time"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Process!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			timedContext, cancel := context.WithTimeout(request.Context(), time.Second)
			defer cancel()

			select {
			case <-timedContext.Done():
				response.ReportError(errors.New("timed out"))
			case <-time.After(time.Minute):
				response.Reply("Processing done!")
			}
		},
	}

	bot.Command("process", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 9

Showcasing the ability to add attachments to a `Reply`

```go
package main

import (
	"context"
	"log"

	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	definition := &slacker.CommandDefinition{
		Description: "Echo a word!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")

			attachments := []slack.Attachment{}
			attachments = append(attachments, slack.Attachment{
				Color:      "red",
				AuthorName: "Raed Shomali",
				Title:      "Attachment Title",
				Text:       "Attachment Text",
			})

			response.Reply(word, slacker.WithAttachments(attachments))
		},
	}

	bot.Command("echo <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 10

Showcasing the ability to create custom responses via `CustomResponse`

```go
package main

import (
	"log"

	"context"
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker"
)

const (
	errorFormat = "> Custom Error: _%s_"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.CustomResponse(NewCustomResponseWriter)

	definition := &slacker.CommandDefinition{
		Description: "Custom!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("custom")
			response.ReportError(errors.New("oops"))
		},
	}

	bot.Command("custom", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// NewCustomResponseWriter creates a new ResponseWriter structure
func NewCustomResponseWriter(channel string, client *slack.Client, rtm *slack.RTM) slacker.ResponseWriter {
	return &MyCustomResponseWriter{channel: channel, client: client, rtm: rtm}
}

// MyCustomResponseWriter a custom response writer
type MyCustomResponseWriter struct {
	channel string
	client  *slack.Client
	rtm     *slack.RTM
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *MyCustomResponseWriter) ReportError(err error) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.channel))
}

// Typing send a typing indicator
func (r *MyCustomResponseWriter) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.channel))
}

// Reply send a attachments to the current channel with a message
func (r *MyCustomResponseWriter) Reply(message string, options ...slacker.ReplyOption) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(message, r.channel))
}

// RTM returns the RTM client
func (r *MyCustomResponseWriter) RTM() *slack.RTM {
	return r.rtm
}

// Client returns the slack client
func (r *MyCustomResponseWriter) Client() *slack.Client {
	return r.client
}
```

## Example 11

Showcasing the ability to toggle the slack Debug option via `WithDebug`

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>", slacker.WithDebug(true))

	definition := &slacker.CommandDefinition{
		Description: "Ping!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	}

	bot.Command("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example 12

Defining a command that can only be executed by authorized users

```go
package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	authorizedUsers := []string{"<USER ID>"}

	authorizedDefinition := &slacker.CommandDefinition{
		Description: "Very secret stuff",
		AuthorizationFunc: func(request slacker.Request) bool {
			return contains(authorizedUsers, request.Event().User)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("You are authorized!")
		},
	}

	bot.Command("secret", authorizedDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func contains(list []string, element string) bool {
	for _, value := range list {
		if value == element {
			return true
		}
	}
	return false
}
```

## Example 13

Adding handlers to when the bot is connected, encounters an error and a default for when none of the commands match

```go
package main

import (
	"log"

	"context"
	"fmt"
	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.Init(func() {
		log.Println("Connected!")
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.DefaultCommand(func(request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what?")
	})

	bot.DefaultEvent(func(event interface{}) {
		fmt.Println(event)
	})

	definition := &slacker.CommandDefinition{
		Description: "help!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("Your own help function...")
		},
	}

	bot.Help(definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```
