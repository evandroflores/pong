package model

import (
	"fmt"

	"github.com/evandroflores/udpong/database"
	"github.com/jinzhu/gorm"
)

func init() {
	database.Connection.AutoMigrate(&Player{})
}

// Player stores the player points and personal info.
type Player struct {
	gorm.Model
	TeamID string `gorm:"not null"`
	ID     string `gorm:"not null"`
	Name   string `gorm:"not null"`
	Image  string
	Points float64
}

// GetPlayer returns a player given ID
func GetPlayer(playerID string) (Player, error) {
	result := Player{}

	database.Connection.Where(&Player{ID: playerID}).
		First(&result)

	return result, nil
}

// Add adds a new player.
func (player Player) Add() {
	database.Connection.Create(&player)
}

// Update a given player
func (player Player) Update() error {
	dbPlayer, err := GetPlayer(player.ID)

	if err != nil {
		return err
	}

	if dbPlayer == (Player{}) {
		return fmt.Errorf("Player not found")
	}

	dbPlayer.Points = player.Points

	database.Connection.Save(&dbPlayer)

	return nil
}
