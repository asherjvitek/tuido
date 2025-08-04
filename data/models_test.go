package data

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPositionWhenEmpty(t *testing.T) {
	expected := defaultPosition
	positionals := []Board{}

	pos, err := GetPosition(positionals, 0)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, expected, pos, fmt.Sprintf("Expected position at index 0 to be %.2f, got %.2f", expected, pos))
}

func TestGetPositionLessThan0(t *testing.T) {
	// expected := 5.0
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
	}

	_, err := GetPosition(positionals, -1)

	assert.EqualError(t, err, "dest should never be less than 0")
}

func TestGetPositionGreaterThanLen(t *testing.T) {
	// expected := 5.0
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
	}

	_, err := GetPosition(positionals, 2)

	assert.EqualError(t, err, "dest should never be > len(positionls) = 1, dest = 2")
}

func TestGetPositionAtBeginning(t *testing.T) {
	expected := 5.0
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
	}

	pos, err := GetPosition(positionals, 0)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, expected, pos, fmt.Sprintf("Expected position at index 0 to be %.2f, got %.2f", expected, pos))
}

func TestGetPositionInbetweenTwo(t *testing.T) {
	expected := 15.0
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
		{BoardId: 2, Name: "Board 2", Position: 20.0},
	}

	pos, err := GetPosition(positionals, 1)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, expected, pos, fmt.Sprintf("Expected position at index 0 to be %.2f, got %.2f", expected, pos))

}

func TestGetPositionAtEnd(t *testing.T) {
	expected := (math.MaxFloat64-10.0)/2 + 10
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
	}

	pos, err := GetPosition(positionals, 1)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, expected, pos, fmt.Sprintf("Expected position at index 0 to be %.2f, got %.2f", expected, pos))

}
