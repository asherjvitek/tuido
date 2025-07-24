package main

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// I think that there must be some better way to do this but it works for the moment.
// I wonder if this would perform if you had like way too many items in one list like I
// do in trello in the completed bucket
func countLines(str string) int {
	lines := 1
	for _, r := range str {
		if r == '\n' {
			lines++
		}
	}

	return lines
}

func getTextInput() textinput.Model {
	input := textinput.New()
	input.Prompt = ""

	return input
}

func getBoards() boards {
	return boards{
		selected: 0,
		boards: []board{
			{
				id:   1,
				name: "Board 1",
				lists: []list{
					{
						title: "TODO",
						items: []string{
							"Thing to Do 6",
						},
					},
					{
						title: "DOING",
						items: []string{
							"Doing this thing 1",
						},
					},
					{
						title: "DONE",
						items: []string{
							"This is done 3!",
						},
					},
				},

				selectedList: 0,
				selectedItem: 0,
				input:        getTextInput(),
			},
			{
				id:   2,
				name: "Board 2",
				lists: []list{
					{
						title: "TODO",
						items: []string{
							"Thing to Do 6",
						},
					},
					{
						title: "DOING",
						items: []string{
							"Doing this thing 1",
						},
					},
					{
						title: "DONE",
						items: []string{
							"This is done 3!",
						},
					},
				},

				selectedList: 0,
				selectedItem: 0,
				input:        getTextInput(),
			},
		},
	}
}
