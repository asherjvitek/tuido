package data

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPositionAt0(t *testing.T) {
	expected := 5.0
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
		{BoardId: 2, Name: "Board 2", Position: 20.0},
	}

	pos := GetPosition(positionals, 0)

	assert.Equal(t, expected, pos, fmt.Sprintf("Expected position at index 0 to be %.2f, got %.2f", expected, pos))

}

func TestGetPositionAt1(t *testing.T) {
	expected := 15.0
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
		{BoardId: 2, Name: "Board 2", Position: 20.0},
	}

	pos := GetPosition(positionals, 1)

	assert.Equal(t, expected, pos, fmt.Sprintf("Expected position at index 0 to be %.2f, got %.2f", expected, pos))

}

func TestGetPositionAtEnd(t *testing.T) {
	expected := math.MaxFloat64 / 2.0 - 20 / 2
	positionals := []Board{
		{BoardId: 1, Name: "Board 1", Position: 10.0},
		{BoardId: 2, Name: "Board 2", Position: 20.0},
	}

	pos := GetPosition(positionals, 3)

	assert.Equal(t, expected, pos, fmt.Sprintf("Expected position at index 0 to be %.2f, got %.2f", expected, pos))

}
