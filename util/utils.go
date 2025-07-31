package util

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"time"
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

func Log(message string) {
	// if len(os.Getenv("DEBUG")) > 0 {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	date := time.Now().Format(time.DateTime)
	fmt.Fprintf(f, "%s: %s", date, message)
	// }
}

func Error(context string, err error) {
	fmt.Printf("%s: %+v", context, err)
	os.Exit(1)
}
