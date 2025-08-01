package data

import (
	"math"
)

type Positional interface {
	GetPosition() float64
}

// TODO: This is not done and I think just needs some tests
func GetPosition[S ~[]E, E Positional](positionals S, dest int) float64 {
	if len(positionals) == 0 {
		return math.MaxFloat64 / 2.0
	}

	if dest == 0 {
		return positionals[dest].GetPosition() / 2
	}

	if dest == len(positionals) {
		return (positionals[dest-1].GetPosition() + math.MaxFloat64) / 2
	}

	return (positionals[dest - 1].GetPosition() + positionals[dest].GetPosition())/2

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
