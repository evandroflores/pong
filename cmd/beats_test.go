package cmd

import (
	"fmt"
	"testing"

	"github.com/nlopes/slack"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/elo"
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
	testILost.loser = makeTestPlayer()
	testILost.evt = makeTestEvent()
	testILost.evt.User = testILost.loser.SlackID

	testIWon := new(BeatsTestSuite)
	testIWon.cmd = iWon
	testIWon.winner = makeTestPlayer()
	testIWon.evt = makeTestEvent()
	testIWon.evt.User = testIWon.winner.SlackID

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

func (s *BeatsTestSuite) TearDownTest() {
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

	blocks := response.GetBlocks()

	s.Len(response.GetBlocks(), 1)
	s.Equal(slack.MBTContext, blocks[0].BlockType())
	contextBlock := blocks[0].(*slack.ContextBlock)

	elements := contextBlock.ContextElements.Elements
	s.Len(elements, 5)

	s.Equal(slack.MixedElementImage, elements[0].MixedElementType())
	s.Equal(s.winner.Image, elements[0].(*slack.ImageBlockElement).ImageURL)
	s.Equal(s.winner.Name, elements[0].(*slack.ImageBlockElement).AltText)

	s.Equal(slack.MixedElementText, elements[1].MixedElementType())
	s.Equal(fmt.Sprintf("*%s* (%04.f pts) *#%02d*",
		s.winner.Name, eloWinnerPts, 1),
		elements[1].(*slack.TextBlockObject).Text)

	s.Equal(slack.MixedElementText, elements[2].MixedElementType())
	s.Equal(" x ", elements[2].(*slack.TextBlockObject).Text)

	s.Equal(slack.MixedElementImage, elements[3].MixedElementType())
	s.Equal(s.loser.Image, elements[3].(*slack.ImageBlockElement).ImageURL)
	s.Equal(s.loser.Name, elements[3].(*slack.ImageBlockElement).AltText)

	s.Equal(slack.MixedElementText, elements[4].MixedElementType())
	s.Equal(fmt.Sprintf("*%s* (%04.f pts) *#%02d*",
		s.loser.Name, eloLoserPts, 2),
		elements[4].(*slack.TextBlockObject).Text)

	s.Empty(response.GetErrors())

}
