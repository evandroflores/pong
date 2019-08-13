package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/suite"
)

type RankTestSuite struct {
	suite.Suite
	rankChannelID   string
	noRankChannelID string
	expectedRank    string
	players         []model.Player
}

func (s *RankTestSuite) SetupSuite() {
	s.rankChannelID = "CCCCCCCCC"
	s.noRankChannelID = "CAAAAAAAA"
	s.expectedRank = fmt.Sprintf("\n*Rank for * <#%s>\n\n", s.rankChannelID)

	s.players = []model.Player{}
	for i := 1; i <= 20; i++ {
		player := makeTestPlayer()
		player.Points = 1000 - float64(i)
		s.players = append(s.players, player)

		database.Connection.Create(&player)
		s.expectedRank += fmt.Sprintf("*%02d* - %04.f - %s\n", i, player.Points, player.Name)
	}
}

func (s *RankTestSuite) TearDownSuite() {
	for i := 0; i < len(s.players); i++ {
		database.Connection.Where(&s.players[i]).Delete(&model.Player{})
	}
	database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})
}

func (s *RankTestSuite) TestMakeEmptyRank() {
	expected := fmt.Sprintf("No rank for channel <#%s>\n\n", s.rankChannelID)
	actual := makeRank(s.rankChannelID, []model.Player{})
	s.Equal(expected, actual)
}

func (s *RankTestSuite) TestMakeRank() {
	actual := makeRank(s.rankChannelID, s.players)
	s.Equal(s.expectedRank, actual)
}

func (s *RankTestSuite) TestEmptyRankForChannel() {
	var props = proper.NewProperties(map[string]string{})

	evt := makeTestEvent()
	evt.Msg.Channel = s.noRankChannelID
	expected := fmt.Sprintf("No rank for channel <#%s>\n\n", s.noRankChannelID)
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	rank(request, response)
	s.Contains(response.GetMessages(), expected)
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}

func (s *RankTestSuite) TestRankForChannel() {
	var props = proper.NewProperties(map[string]string{})

	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	rank(request, response)
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
	s.Equal(s.expectedRank, response.GetMessages()[0])
}

func TestRankTestSuite(t *testing.T) {
	suite.Run(t, new(RankTestSuite))
}
