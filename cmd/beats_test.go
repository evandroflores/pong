package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/elo"
	"github.com/nlopes/slack"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/shomali11/slacker"
	"github.com/stretchr/testify/suite"
)

type BeatsTestSuite struct {
	suite.Suite
	winner               model.Player
	originalWinnerPoints float64
	loser                model.Player
	originalLoserPoints  float64
	evt                  *slack.MessageEvent
	cmd                  func(request slacker.Request, response slacker.ResponseWriter)
}

func TestBeatsTestSuite(t *testing.T) {
	testBeats := new(BeatsTestSuite)
	testBeats.cmd = beats
	testBeats.evt = makeTestEvent()

	testILost := new(BeatsTestSuite)
	testILost.cmd = iLost
	testILost.winner = makeTestPlayer()
	testILost.evt = makeTestEvent()
	testILost.evt.User = testILost.winner.SlackID

	testIWon := new(BeatsTestSuite)
	testIWon.cmd = iWon
	testIWon.loser = makeTestPlayer()
	testIWon.evt = makeTestEvent()
	testIWon.evt.User = testILost.loser.SlackID

	suite.Run(t, testBeats)
	suite.Run(t, testILost)
	suite.Run(t, testIWon)
}

func (s *BeatsTestSuite) SetupTest() {
	s.originalWinnerPoints = 1200
	s.originalLoserPoints = 800
	if (s.winner == model.Player{}) {
		s.winner = makeTestPlayer()
	}
	s.winner.Points = s.originalWinnerPoints
	if (s.loser == model.Player{}) {
		s.loser = makeTestPlayer()
	}
	s.loser.Points = s.originalLoserPoints

	database.Connection.Create(&s.winner)
	database.Connection.Create(&s.loser)
}

func (s *BeatsTestSuite) TearDownSuite() {
	database.Connection.Unscoped().Delete(&s.winner)
	database.Connection.Unscoped().Delete(&s.loser)
}

func (s *BeatsTestSuite) TestExpectedEloResult() {
	var props = proper.NewProperties(
		map[string]string{
			"@winner": s.winner.SlackID,
			"@loser":  s.loser.SlackID,
		})

	request := &fakeRequest{event: s.evt, properties: props}
	response := &fakeResponse{}

	eloWinnerPts, eloLoserPts := elo.Calc(s.originalWinnerPoints, s.originalLoserPoints)
	s.cmd(request, response)

	s.Contains(response.GetMessages(),
		fmt.Sprintf("*%s* %04.f pts (#%02d) vs *%s* %04.f pts (#%02d)",
			s.winner.Name, eloWinnerPts, 1, s.loser.Name, eloLoserPts, 2))

	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}
