package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/elo"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/suite"
)

type SimulateTestSuite struct {
	suite.Suite
	playerA               model.Player
	originalPlayerAPoints float64
	playerB               model.Player
	originalPlayerBPoints float64
}

func TestSimulateTestSuite(t *testing.T) {
	suite.Run(t, new(SimulateTestSuite))
}

func (s *SimulateTestSuite) SetupSuite() {
	s.originalPlayerAPoints = 1200
	s.originalPlayerBPoints = 800
	s.playerA = makeTestPlayer()
	s.playerA.Points = s.originalPlayerAPoints
	s.playerB = makeTestPlayer()
	s.playerB.Points = s.originalPlayerBPoints

	database.Connection.Create(&s.playerA)
	database.Connection.Create(&s.playerB)
}

func (s *SimulateTestSuite) TestWinnerNotAUser() {
	var props = proper.NewProperties(
		map[string]string{
			"@playerA": s.playerA.Name,
			"@playerB": s.playerB.SlackID,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	simulate(request, response)
	s.Contains(response.GetMessages(), "_The given winner is not a User_")
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}

func (s *SimulateTestSuite) TestLoserNotAUser() {
	var props = proper.NewProperties(
		map[string]string{
			"@playerA": s.playerA.SlackID,
			"@playerB": s.playerB.Name,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	simulate(request, response)
	s.Contains(response.GetMessages(), "_The given loser is not a User_")
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}

func (s *SimulateTestSuite) TestExpectedEloResult() {
	var props = proper.NewProperties(
		map[string]string{
			"@playerA": s.playerA.SlackID,
			"@playerB": s.playerB.SlackID,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	eloAWin, eloBLose := elo.Calc(s.originalPlayerAPoints, s.originalPlayerBPoints)
	eloBWin, eloALose := elo.Calc(s.originalPlayerBPoints, s.originalPlayerAPoints)

	simulate(request, response)

	s.Contains(response.GetMessages(),
		fmt.Sprintf("*%s* wins (%04.f pts) vs *%s* (%04.f pts)",
			s.playerA.Name, eloAWin, s.playerB.Name, eloBLose))

	s.Contains(response.GetMessages(),
		fmt.Sprintf("*%s* wins (%04.f pts) vs *%s* (%04.f pts)",
			s.playerB.Name, eloBWin, s.playerA.Name, eloALose))

	s.Len(response.GetMessages(), 2)
	s.Empty(response.GetErrors())
}

func (s *SimulateTestSuite) TestMakeSureSimulateKeepsPointsUnchanged() {
	var props = proper.NewProperties(
		map[string]string{
			"@playerA": s.playerA.SlackID,
			"@playerB": s.playerB.SlackID,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	simulate(request, response)
	s.Empty(response.GetErrors())

	playerAFromDB, errA := model.GetPlayer(s.playerA.TeamID, s.playerA.ChannelID, s.playerA.SlackID)
	playerBFromDB, errB := model.GetPlayer(s.playerB.TeamID, s.playerB.ChannelID, s.playerB.SlackID)

	s.NoError(errA)
	s.NoError(errB)

	s.Equal(s.originalPlayerAPoints, playerAFromDB.Points)
	s.Equal(s.originalPlayerBPoints, playerBFromDB.Points)
}

func (s *SimulateTestSuite) TearDownSuite() {
	database.Connection.Where(&s.playerA).Or(&s.playerB).Delete(&model.Player{})
	database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})
}
