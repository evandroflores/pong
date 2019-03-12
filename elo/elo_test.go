package elo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartersCalc(t *testing.T) {
	expectWinner := 1016.0
	expectLoser := 984.0
	updatedWinner, updatedLoser := Calc(1000, 1000)

	assert.Equal(t, expectWinner, updatedWinner)
	assert.Equal(t, expectLoser, updatedLoser)
}

func TestBigLeapCalc(t *testing.T) {
	expectWinner := 829.0
	expectLoser := 1171.0
	updatedWinner, updatedLoser := Calc(800, 1200)

	assert.Equal(t, math.Round(expectWinner), math.Round(updatedWinner))
	assert.Equal(t, math.Round(expectLoser), math.Round(updatedLoser))
}

func TestPrecisionCalc(t *testing.T) {
	expectWinner := 997.8159093133837
	expectLoser := 997.8130906866162
	updatedWinner, updatedLoser := Calc(980.199, 1015.43)

	assert.Equal(t, expectWinner, updatedWinner)
	assert.Equal(t, expectLoser, updatedLoser)
}
