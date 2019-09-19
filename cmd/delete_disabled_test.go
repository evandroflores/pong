package cmd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/evandroflores/pong/slack"
	ns "github.com/nlopes/slack"
	"github.com/shomali11/proper"
	"github.com/shomali11/slacker"
	"github.com/stretchr/testify/suite"
)

type DeleteDisabledTestSuite struct {
	suite.Suite
	teamID      string
	channelID   string
	adminID     string
	nonAdminID  string
	players     []model.Player
	deactivated map[string]model.Player
}

func (s *DeleteDisabledTestSuite) SetupSuite() {
	s.teamID = "TDELETE"
	s.channelID = "CTESTDEL"
	s.adminID = "UIS0ADMIN"
	s.nonAdminID = "UNOTADMIN"
	s.players = []model.Player{}
	s.deactivated = map[string]model.Player{}

	for i := 1; i <= 20; i++ {
		player := makeTestPlayerWith(s.teamID, s.channelID, fmt.Sprintf("U%08d", i), fmt.Sprintf("Fake User %08d", i))
		player.Points = 1000 - float64(i)
		if i%2 == 0 {
			s.deactivated[player.SlackID] = player
		}
		s.players = append(s.players, player)

		database.Connection.Create(&player)
	}

	mockGetUserInfo := func(sl *slacker.Slacker, id string) (*ns.User, error) {
		isDeactivated := false
		if (s.deactivated[id] != model.Player{}) {
			isDeactivated = true
		}

		return &ns.User{
			ID:      id,
			Deleted: isDeactivated,
		}, nil
	}

	monkey.PatchInstanceMethod(reflect.TypeOf(slack.Client), "GetUserInfo", mockGetUserInfo)
}

func (s *DeleteDisabledTestSuite) TearDownSuite() {
	for i := 0; i < len(s.players); i++ {
		database.Connection.Where(&s.players[i]).Delete(&model.Player{})
	}
	database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})
}

func TestDeleteDisabledTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteDisabledTestSuite))
}

func (s *DeleteDisabledTestSuite) TestAdminRemoving() {
	var props = proper.NewProperties(map[string]string{})

	evt := makeTestEventWith(s.teamID, s.channelID, s.adminID)

	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	deleteDisabled(request, response)
	s.Len(response.GetMessages(), 1)
	s.Equal(strings.Count(response.GetMessages()[0], "Removed"), len(s.deactivated))
	s.Empty(response.GetErrors())
}

func (s *DeleteDisabledTestSuite) TestNoUsersRemoved() {
	var props = proper.NewProperties(map[string]string{})

	evt := makeTestEventWith(s.teamID, s.channelID, s.adminID)

	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	deleteDisabled(request, response)
	s.Len(response.GetMessages(), 1)
	s.Equal(response.GetMessages()[0], "No users removed.")
	s.Empty(response.GetErrors())
}
