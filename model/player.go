package model

import (
	"fmt"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/slack"
	"github.com/jinzhu/gorm"
	ns "github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

func init() {
	database.Connection.AutoMigrate(&Player{})
}

// Player stores the player points and personal info.
type Player struct {
	gorm.Model
	TeamID    string `gorm:"not null"`
	ChannelID string `gorm:"not null"`
	SlackID   string `gorm:"not null"`
	Name      string `gorm:"not null"`
	Image     string
	Points    float64 `gorm:"default:1000"`
}

// ToStr returns a string representation of Player
func (player *Player) ToStr() string {
	return fmt.Sprintf("TeamID: %s ChannelID: %s SlackID: %s", player.TeamID, player.ChannelID, player.SlackID)
}

// IDStr returns a simplified representation of Player IDs
func (player *Player) IDStr() string {
	return fmt.Sprintf("%s.%s.%s", player.TeamID, player.ChannelID, player.SlackID)
}

// GetPlayer returns a player given TeamID+ChannelID+SlackID
func GetPlayer(teamID, channelID, slackID string) (Player, error) {
	result := Player{}

	database.Connection.Where(&Player{TeamID: teamID, ChannelID: channelID, SlackID: slackID}).
		First(&result)

	return result, nil
}

// Add adds a new player.
func (player *Player) Add() {
	player.IngestData()
	database.Connection.Create(&player)
}

// GetOrCreatePlayer will try to find the player, if can't, will return an brand new one
func GetOrCreatePlayer(teamID, channelID, slackID string) (Player, error) {
	player, err := GetPlayer(teamID, channelID, slackID)
	if err != nil {
		return Player{}, err
	}
	if player == (Player{}) {
		tempPlayer := Player{TeamID: teamID, ChannelID: channelID, SlackID: slackID}
		log.Debugf("Player not found, will create... %s", tempPlayer.ToStr())
		tempPlayer.Add()
		return GetPlayer(teamID, channelID, slackID)
	}
	return player, nil
}

// Update a given player
func (player *Player) Update() (int, error) {
	dbPlayer, err := GetPlayer(player.TeamID, player.ChannelID, player.SlackID)

	if err != nil {
		return 0, err
	}

	if dbPlayer == (Player{}) {
		return 0, fmt.Errorf("player not found - %s", player.ToStr())
	}

	previousPosition := dbPlayer.GetPosition()
	dbPlayer.IngestData()
	dbPlayer.Points = player.Points

	database.Connection.Save(&dbPlayer)

	return (previousPosition - dbPlayer.GetPosition()), nil
}

// Delete removes locally the user from the DB.
func (player *Player) Delete() {
	database.Connection.Delete(&player)
}

// IngestData calls Slack API to get Users data
func (player *Player) IngestData() {
	slackUser, err := slack.Client.GetUserInfo(player.SlackID)
	if err != nil {
		log.Errorf("Something bad happen while ingesting User data from Slack - %s", err.Error())
		return
	}
	player.Name = slackUser.RealName
	player.Image = slackUser.Profile.Image48
	player.TeamID = slackUser.TeamID
}

// GetPosition returns the current position for the given player
func (player *Player) GetPosition() int {
	count := 0

	database.Connection.
		Model(&Player{}).
		Where(&Player{TeamID: player.TeamID, ChannelID: player.ChannelID}).
		Where("points > ?", player.Points).
		Count(&count)

	return count + 1
}

// GetPlayers list a number of players limited by the given threshold
func GetPlayers(teamID, channelID string, limit int) []Player {
	results := []Player{}

	database.Connection.
		Where(&Player{TeamID: teamID, ChannelID: channelID}).
		Order("points desc, created_at").
		Limit(limit).
		Find(&results)

	return results
}

// GetAllPlayers list all players
func GetAllPlayers(teamID, channelID string) []Player {
	results := []Player{}

	database.Connection.
		Where(&Player{TeamID: teamID, ChannelID: channelID}).
		Order("points desc, created_at").
		Find(&results)

	return results
}

// GetBlockCardWithText formats the Player Image and Name on a Slack Block Kit
func (player *Player) GetBlockCardWithText(before, after string) []ns.MixedElement {
	elements := []ns.MixedElement{}

	if before != "" {
		beforeTextBlock := ns.NewTextBlockObject(ns.MarkdownType, before, false, false)
		elements = append(elements, beforeTextBlock)
	}
	avatar := ns.NewImageBlockElement(player.Image, player.Name)
	elements = append(elements, avatar)

	name := ns.NewTextBlockObject(ns.MarkdownType, fmt.Sprintf("*%s*", player.Name), false, false)
	elements = append(elements, name)

	if after != "" {
		afterTextBlock := ns.NewTextBlockObject(ns.MarkdownType, after, false, false)
		elements = append(elements, afterTextBlock)
	}

	return elements
}

// GetBlockCard calls GetBlockCardWithText using a default text
func (player *Player) GetBlockCard() []ns.MixedElement {
	return player.GetBlockCardWithText(
		fmt.Sprintf("*#%02d*", player.GetPosition()),
		fmt.Sprintf("(%04.f pts)", player.Points))
}
