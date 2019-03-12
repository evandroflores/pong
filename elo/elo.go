package elo

import "math"

// https://en.wikipedia.org/wiki/Elo_rating_system
const (
	k         float64 = 32
	deviation float64 = 400
)

// Calc returns the updated Winner and Loser
func Calc(winner, loser float64) (updatedWinner, updatedLoser float64) {
	prob := 1 / (1 + math.Pow(10, (loser-winner)/deviation))
	diff := k * (1 - prob)

	updatedWinner = winner + diff
	updatedLoser = loser - diff

	return updatedWinner, updatedLoser
}
