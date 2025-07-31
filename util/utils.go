package util

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// I think that there must be some better way to do this but it works for the moment.
// I wonder if this would perform if you had like way too many items in one list like I
// do in trello in the completed bucket
func CountLines(str string) int {
	lines := 1
	for _, r := range str {
		if r == '\n' {
			lines++
		}
	}

	return lines
}

func GetTextInput() textinput.Model {
	input := textinput.New()
	input.Prompt = ""

	return input
}

// func getDefaultBoards() boards.Model {
// 	return boards.Model{
// 		Boards: []board.Model{
// 			{
// 				Id:   1,
// 				Name: "My First Board",
// 				Lists: []board.List{
// 					{
// 						Title: "TODO",
// 						Items: []string{
// 							"Thing to Do 6",
// 						},
// 					},
// 				},
//
// 				selectedList: 0,
// 				selectedItem: 0,
// 				input:        getTextInput(),
// 			},
// 		},
// 	}
// }
