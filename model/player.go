package model

import (
	"fmt"

	"github.com/evandroflores/pong/database"
	"github.com/evandroflores/pong/slack"
	"github.com/jinzhu/gorm"
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
func (player Player) ToStr() string {
	return fmt.Sprintf("TeamID: %s ChannelID: %s SlackID: %s", player.TeamID, player.ChannelID, player.SlackID)
}

// GetPlayer returns a player given TeamID+ChannelID+SlackID
func GetPlayer(teamID, channelID, slackID string) (Player, error) {
	result := Player{}

	database.Connection.Where(&Player{TeamID: teamID, ChannelID: channelID, SlackID: slackID}).
		First(&result)

	return result, nil
}

// Add adds a new player.
func (player Player) Add() {
	player.ingestData()
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
func (player Player) Update() error {
	dbPlayer, err := GetPlayer(player.TeamID, player.ChannelID, player.SlackID)

	if err != nil {
		return err
	}

	if dbPlayer == (Player{}) {
		return fmt.Errorf("Player not found - %s", player.ToStr())
	}

	dbPlayer.ingestData()
	dbPlayer.Points = player.Points

	database.Connection.Save(&dbPlayer)

	return nil
}

func (player *Player) ingestData() {
	slackUser, err := slack.Client.GetUserInfo(player.SlackID)
	if err != nil {
		log.Errorf("Something bad happen while ingesting User data from Slack - %s", err.Error())
		return
	}
	player.Name = slackUser.RealName
	player.Image = slackUser.Profile.Image48
	player.TeamID = slackUser.TeamID
}

// GetPlayers list a number of players limited by the given threshhold
func GetPlayers(teamID string, channelID string, limit int) []Player {
	results := []Player{}

	database.Connection.Find(&results).
		Where(&Player{TeamID: teamID, ChannelID: channelID}).
		Order("points, created_at").
		Limit(limit)

	return results
}

// GetAllPlayers list all players
func GetAllPlayers(teamID string, channelID string) []Player {
	results := []Player{}

	database.Connection.Find(&results).
		Where(&Player{TeamID: teamID, ChannelID: channelID}).
		Order("points, created_at")

	return results
}
