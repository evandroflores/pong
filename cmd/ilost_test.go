package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/elo"
	"github.com/nlopes/slack"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/suite"
)

type LostTestSuite struct {
	suite.Suite
	winner               model.Player
	originalWinnerPoints float64
	loser                model.Player
	originalLoserPoints  float64
	evt                  *slack.MessageEvent
}

func TestLostTestSuite(t *testing.T) {
	suite.Run(t, new(LostTestSuite))
}

func (s *LostTestSuite) SetupSuite() {
	s.originalWinnerPoints = 1200
	s.originalLoserPoints = 800
	s.winner = makeTestPlayer()
	s.winner.Points = s.originalWinnerPoints
	s.loser = makeTestPlayer()
	s.loser.Points = s.originalLoserPoints
	database.Connection.Create(&s.winner)
	database.Connection.Create(&s.loser)

	s.evt = makeTestEvent()
	s.evt.User = s.loser.SlackID
}

func (s *LostTestSuite) TearDownSuite() {
	database.Connection.Unscoped().Delete(&s.winner)
	database.Connection.Unscoped().Delete(&s.loser)
}

func (s *LostTestSuite) TestExpectedEloResult() {
	var props = proper.NewProperties(
		map[string]string{
			"@winner": s.winner.SlackID,
		})

	request := &fakeRequest{event: s.evt, properties: props}
	response := &fakeResponse{}

	eloWinnerPts, eloLoserPts := elo.Calc(s.originalWinnerPoints, s.originalLoserPoints)

	iLost(request, response)

	s.Contains(response.GetMessages(),
		fmt.Sprintf("*%s* %04.f pts (#%02d) vs *%s* %04.f pts (#%02d)",
			s.winner.Name, eloWinnerPts, 1, s.loser.Name, eloLoserPts, 2))

	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}
