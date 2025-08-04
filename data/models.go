package data

import (
	"fmt"
	"math"
)

type Positional interface {
	GetPosition() float64
}

const defaultPosition = math.MaxFloat64 / 2.0

func GetPosition[S ~[]E, E Positional](positionals S, dest int) (float64, error) {
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
		return (math.MaxFloat64 - pos) / 2 + pos, nil
	}

	a := positionals[dest].GetPosition()
	b := positionals[dest - 1].GetPosition()

	return (a - b) / 2 + b, nil
}

type Board struct {
	BoardId  int
	Name     string
	Position float64
}

func (b Board) GetPosition() float64 {
	return b.Position
}

type List struct {
	ListId   int
	BoardId  int
	Name     string
	Position float64
	Items    []Item
}

func (b List) GetPosition() float64 {
	return b.Position
}

type Item struct {
	ItemId   int
	ListId   int
	Text     string
	Position float64
}

func (b Item) GetPosition() float64 {
	return b.Position
}
