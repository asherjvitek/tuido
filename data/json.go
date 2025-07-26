package data

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"errors"
	"path/filepath"
	"tuido/boards"
)

// This is the first pass at something comepletely simple and naive but it should get persistence and loading.
const name = "boards.json"

func getDir() string {
	user, err := user.Current()

	if err != nil {
		return os.TempDir()
	}

	return user.HomeDir
}

func SaveData(boards boards.Model) {

	jsonData, err := json.Marshal(boards)

	if err != nil {
		panic("Something broke with the jsonData")
	}

	dir := getDir()
	path := filepath.Join(dir, name)

	err = os.WriteFile(path, jsonData, 0644)

	if err != nil {
		panic(fmt.Sprintf("Failed to Write file %s", path))
	}
}

func LoadData() boards.Model {
	var boards boards.Model
	dir := getDir()
	path := filepath.Join(dir, name)
	data, err := os.ReadFile(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return boards.New()
		}

		panic(fmt.Sprintf("There was some other error trying to read the file %s\nerror: %s", path, err.Error()))
	}

	err = json.Unmarshal(data, &boards)

	if err != nil {
		panic(fmt.Sprintf("Failed to Deserialize data from %s", path))
	}

	//We have to do this to make sure that everything is setup right
	for i := range boards.Boards {
		boards.Boards[i].Setup()
	}

	return boards

}
