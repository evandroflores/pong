package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/suite"
)

type TopTestSuite struct {
	suite.Suite
	rankChannelID   string
	noRankChannelID string
	top20rank       string
	top10rank       string
	top5rank        string
	players         []model.Player
}

func (s *TopTestSuite) SetupSuite() {
	s.rankChannelID = "CCCCCCCCC"
	s.noRankChannelID = "CAAAAAAAA"
	rankHeader := fmt.Sprintf("\n*Rank for * <#%s>\n\n", s.rankChannelID)
	s.top20rank = rankHeader
	s.top10rank = rankHeader
	s.top5rank = rankHeader

	s.players = []model.Player{}
	for i := 1; i <= 20; i++ {
		player := makeTestPlayer()
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

func (s *TopTestSuite) TearDownSuite() {
	for i := 0; i < len(s.players); i++ {
		database.Connection.Where(&s.players[i]).Delete(&model.Player{})
	}
	database.Connection.Unscoped().Where("deleted_at is not null").Delete(&model.Player{})
}

func (s *TopTestSuite) TestEmptyTop() {
	var props = proper.NewProperties(map[string]string{})

	evt := makeTestEvent()
	evt.Msg.Channel = s.noRankChannelID
	expected := fmt.Sprintf("No rank for channel <#%s>\n\n", s.noRankChannelID)
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Contains(response.GetMessages(), expected)
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
}

func (s *TopTestSuite) TestTop10WithoutSendingParam() {
	var props = proper.NewProperties(
		map[string]string{})
	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Len(response.GetMessages(), 1)
	s.Empty(response.GetErrors())
	s.Equal(s.top10rank, response.GetMessages()[0])
}

func (s *TopTestSuite) TestTop10() {
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

func (s *TopTestSuite) TestTop5ToGuaranteeIsNotFollowingDefault() {
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

func TestTopTestSuite(t *testing.T) {
	suite.Run(t, new(TopTestSuite))
}
