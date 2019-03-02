package model

import (
	"fmt"

	"github.com/evandroflores/udpong/database"
	"github.com/evandroflores/udpong/slack"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func init() {
	database.Connection.AutoMigrate(&Player{})
}

// Player stores the player points and personal info.
type Player struct {
	gorm.Model
	SlackID string `gorm:"not null"`
	TeamID  string `gorm:"not null"`
	Name    string `gorm:"not null"`
	Image   string
	Points  float64 `gorm:"default:1000"`
}

// GetPlayer returns a player given ID
func GetPlayer(slackID string) (Player, error) {
	result := Player{}

	database.Connection.Where(&Player{SlackID: slackID}).
		First(&result)

	return result, nil
}

// Add adds a new player.
func (player Player) Add() {
	player.ingestData()
	fmt.Printf("Add:%s", player.Name)
	database.Connection.Create(&player)
}

// GetOrCreatePlayer will try to find the player, if can't, will return an brand new one
func GetOrCreatePlayer(slackID string) (Player, error) {
	player, err := GetPlayer(slackID)
	if err != nil {
		return Player{}, err
	}
	if player == (Player{}) {
		Player{SlackID: slackID}.Add()
	}
	return player, nil
}

// Update a given player
func (player Player) Update() error {
	dbPlayer, err := GetPlayer(player.SlackID)

	if err != nil {
		return err
	}

	if dbPlayer == (Player{}) {
		return fmt.Errorf("Player with the ID %s not found", player.SlackID)
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
}
