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
	rankHeader      string
	top20rank       string
	top10rank       string
	top5rank        string
	players         []model.Player
}

func (s *RankTestSuite) SetupSuite() {
	s.rankChannelID = "CCCCCCCCC"
	s.noRankChannelID = "CAAAAAAAA"
	rankHeader := fmt.Sprintf("\n*Rank for * <#%s>\n\n", s.rankChannelID)
	s.top20rank = rankHeader
	s.top10rank = rankHeader
	s.top5rank = rankHeader

	s.players = []model.Player{}
	for i := 1; i <= 20; i++ {
		player := makeTestPlayer()
		player.Name = fmt.Sprintf("%s - %02d", player.Name, i)
		player.SlackID = fmt.Sprintf("UUUUUUU%02d", i)
		player.Points = 1000 - float64(i)
		s.players = append(s.players, player)

		database.Connection.Create(&player)
		s.top20rank += fmt.Sprintf("*%02d* - %04.f - %s\n", i, player.Points, player.Name)
		if i <= 10 {
			s.top10rank += fmt.Sprintf("*%02d* - %04.f - %s\n", i, player.Points, player.Name)
		}
		if i <= 5 {
			s.top5rank += fmt.Sprintf("*%02d* - %04.f - %s\n", i, player.Points, player.Name)
		}
	}
}

func (s *RankTestSuite) TearDownSuite() {
	for i := 0; i < len(s.players); i++ {
		database.Connection.Where(&s.players[i]).Delete(&model.Player{})
	}
	database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})
}

func (s *RankTestSuite) TestMakeEmptyRank() {
	expected := "No rank for this channel"
	actual := makeRank(s.rankChannelID, []model.Player{})
	s.Equal(expected, actual)
}

func (s *RankTestSuite) TestMakeRank() {
	actual := makeRank(s.rankChannelID, s.players)
	s.Equal(s.top20rank, actual)
}

func (s *RankTestSuite) TestEmptyRankForChannel() {
	var props = proper.NewProperties(map[string]string{})

	evt := makeTestEvent()
	evt.Msg.Channel = s.noRankChannelID
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	rank(request, response)
	s.Contains(response.GetMessages(), "No rank for this channel")
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
	s.Equal(s.top20rank, response.GetMessages()[0])
}

// Top Tests *********************************************************************************
func (s *RankTestSuite) TestEmptyTop() {
	var props = proper.NewProperties(map[string]string{})

	evt := makeTestEvent()
	evt.Msg.Channel = s.noRankChannelID
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Contains(response.GetMessages(), "No rank for this channel")
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}

func (s *RankTestSuite) TestTop10WithoutSendingParam() {
	var props = proper.NewProperties(
		map[string]string{})
	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
	s.Equal(s.top10rank, response.GetMessages()[0])
}

func (s *RankTestSuite) TestTop10() {
	var props = proper.NewProperties(
		map[string]string{
			"limit": "10",
		})
	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
	s.Equal(s.top10rank, response.GetMessages()[0])
}

func (s *RankTestSuite) TestTop5ToGuaranteeIsNotFollowingDefault() {
	var props = proper.NewProperties(
		map[string]string{
			"limit": "5",
		})
	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
	s.Equal(s.top5rank, response.GetMessages()[0])
}

func TestRankTestSuite(t *testing.T) {
	suite.Run(t, new(RankTestSuite))
}
