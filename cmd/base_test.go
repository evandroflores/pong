package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const channelIDUnclean string = "<#C12345678>"
const userIDUnclean string = "<@U12345678>"
const teamID string = "T12345678"
const channelID string = "C12345678"
const userID string = "U12345678"

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
