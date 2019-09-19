package cmd

import (
	"fmt"
	"testing"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/model"
	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
	"github.com/stretchr/testify/suite"
)

type TopTestSuite struct {
	suite.Suite
	rankChannelID   string
	noRankChannelID string
	rankHeader      string
	players         []model.Player
}

func (s *TopTestSuite) SetupSuite() {
	s.rankChannelID = "C12345678"
	s.noRankChannelID = "CAAAAAAAA"
	s.rankHeader = fmt.Sprintf("*Rank for * <#%s>", s.rankChannelID)

	s.players = []model.Player{}
	for i := 1; i <= 20; i++ {
		player := makeTestPlayer()
		player.Points = 1000 - float64(i)
		s.players = append(s.players, player)

		database.Connection.Create(&player)
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
	expected := fmt.Sprintf("No rank for channel <#%s>", s.noRankChannelID)
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Contains(response.GetMessages(), expected)
	fmt.Println(response.GetMessages())
	s.Len(response.GetMessages(), 1)
	s.Len(response.GetBlocks(), 0)
	s.Empty(response.GetErrors())
}

func (s *TopTestSuite) TestMoreThanTop20() {
	var props = proper.NewProperties(
		map[string]string{
			"limit": "40",
		})

	evt := makeTestEvent()
	expected := "Top is limited to 20 players"
	request := &fakeRequest{event: evt, properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Contains(response.GetMessages(), expected)
	fmt.Println(response.GetMessages())
	s.Len(response.GetMessages(), 1)
	s.Len(response.GetBlocks(), 0)
	s.Empty(response.GetErrors())
}

func (s *TopTestSuite) TestTop10WithoutSendingParam() {
	var props = proper.NewProperties(
		map[string]string{})
	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Equal(response.GetMessages(), []string{""})
	blocks := response.GetBlocks()
	s.Len(blocks, 11)
	s.Empty(response.GetErrors())

	headerBlock := blocks[0].(*slack.ContextBlock).ContextElements.Elements[0]

	s.Equal(s.rankHeader, headerBlock.(*slack.TextBlockObject).Text)

	for pos := 1; pos <= 10; pos++ {
		elements := blocks[pos].(*slack.ContextBlock).ContextElements.Elements
		player := s.players[pos-1]

		s.Equal(slack.MixedElementText, elements[0].MixedElementType())
		s.Equal(fmt.Sprintf("*%02d* - %04.f", pos, player.Points),
			elements[0].(*slack.TextBlockObject).Text)

		s.Equal(slack.MixedElementImage, elements[1].MixedElementType())
		s.Equal(player.Image, elements[1].(*slack.ImageBlockElement).ImageURL)
		s.Equal(player.Name, elements[1].(*slack.ImageBlockElement).AltText)

		s.Equal(slack.MixedElementText, elements[2].MixedElementType())
		s.Equal(fmt.Sprintf("*%s*", player.Name), elements[2].(*slack.TextBlockObject).Text)
	}
}

func (s *TopTestSuite) TestTop10() {
	var props = proper.NewProperties(
		map[string]string{
			"limit": "10",
		})
	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Equal(response.GetMessages(), []string{""})
	blocks := response.GetBlocks()
	s.Len(blocks, 11)
	s.Empty(response.GetErrors())
}

func (s *TopTestSuite) TestTop5ToGuaranteeIsNotFollowingDefault() {
	var props = proper.NewProperties(
		map[string]string{
			"limit": "5",
		})
	request := &fakeRequest{event: makeTestEvent(), properties: props}
	response := &fakeResponse{}

	top(request, response)
	s.Equal(response.GetMessages(), []string{""})
	blocks := response.GetBlocks()
	s.Len(blocks, 6)
	s.Empty(response.GetErrors())
}

func TestTopTestSuite(t *testing.T) {
	suite.Run(t, new(TopTestSuite))
}
