package util

import (
	"fmt"
	"os/user"
	"path/filepath"

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

func GetAppDir() (string, error) {
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Unable to load the current user, err: %s", err.Error())
		return "", err
	}

	return filepath.Join(user.HomeDir, ".tuido"), nil
}
