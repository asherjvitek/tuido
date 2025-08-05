package data

import (
	"fmt"
	"math"
)

type Positional interface {
	GetPosition() int
}

const defaultPosition = math.MaxInt / 2

func GetPosition[S ~[]E, E Positional](positionals S, dest int) (int, error) {
	if dest < 0 {
		return 0, fmt.Errorf("dest should never be less than 0")
	}

	if dest > len(positionals) {
		return 0, fmt.Errorf("dest should never be > len(positionls) = %d, dest = %d", len(positionals), dest)
	}

	if len(positionals) == 0 {
		return defaultPosition, nil
	}

	if dest == 0 {
		return positionals[dest].GetPosition() / 2, nil
	}

	if dest == len(positionals) {
		pos := positionals[dest-1].GetPosition()
		return (math.MaxInt64 - pos) / 2 + pos, nil
	}

	a := positionals[dest].GetPosition()
	b := positionals[dest - 1].GetPosition()

	return (a - b) / 2 + b, nil
}

type Board struct {
	BoardId  int
	Name     string
	Position int
}

func (b Board) GetPosition() int {
	return b.Position
}

type List struct {
	ListId   int
	BoardId  int
	Name     string
	Position int
	Items    []Item
}

func (b List) GetPosition() int {
	return b.Position
}

type Item struct {
	ItemId   int
	ListId   int
	Text     string
	Position int
}

func (b Item) GetPosition() int {
	return b.Position
}
