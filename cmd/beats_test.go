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

type BeatsTestSuite struct {
	suite.Suite
	winner               model.Player
	originalWinnerPoints float64
	loser                model.Player
	originalLoserPoints  float64
}

func TestBeatsTestSuite(t *testing.T) {
	suite.Run(t, new(BeatsTestSuite))
}

func (s *BeatsTestSuite) SetupSuite() {
	s.originalWinnerPoints = 1200
	s.originalLoserPoints = 800
	s.winner = makeTestPlayer()
	s.winner.Points = s.originalWinnerPoints
	s.loser = makeTestPlayer()
	s.loser.Points = s.originalLoserPoints

	database.Connection.Create(&s.winner)
	database.Connection.Create(&s.loser)
}

func (s *BeatsTestSuite) TearDownSuite() {
	database.Connection.Unscoped().Delete(&s.winner)
	database.Connection.Unscoped().Delete(&s.loser)
}

func (s *BeatsTestSuite) TestWinnerNotAUser() {
	var props = proper.NewProperties(
		map[string]string{
			"@winner": s.winner.Name,
			"@loser":  s.loser.SlackID,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	beats(request, response)
	s.Contains(response.GetErrors(), "the given winner is not a user")
	s.Len(response.GetErrors(), 1)
}

func (s *BeatsTestSuite) TestLoserNotAUser() {
	var props = proper.NewProperties(
		map[string]string{
			"@winner": s.winner.SlackID,
			"@loser":  s.loser.Name,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	beats(request, response)
	s.Contains(response.GetErrors(), "the given loser is not a user")
	s.Len(response.GetErrors(), 1)
}

func (s *BeatsTestSuite) TestExpectedEloResult() {
	var props = proper.NewProperties(
		map[string]string{
			"@winner": s.winner.SlackID,
			"@loser":  s.loser.SlackID,
		})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	eloWinnerPts, eloLoserPts := elo.Calc(s.originalWinnerPoints, s.originalLoserPoints)

	beats(request, response)

	s.Contains(response.GetMessages(),
		fmt.Sprintf("*%s* %04.f pts (#%02d) vs *%s* %04.f pts (#%02d)",
			s.winner.Name, eloWinnerPts, 1, s.loser.Name, eloLoserPts, 2))

	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}
