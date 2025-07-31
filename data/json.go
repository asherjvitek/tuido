package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"tuido/boards"
)

// This is the first pass at something comepletely simple and naive but it should get persistence and loading.
const saveDir = ".tuido"
const fileName = "boards.json"

func getPath() string {
	user, err := user.Current()

	if err != nil {
		fmt.Printf("Unable to loac the current user, err: %s", err.Error())
	}

	return filepath.Join(user.HomeDir, saveDir, fileName)
}

func SaveData(boards boards.Model) {
	jsonData, err := json.Marshal(boards)

	if err != nil {
		panic("Something broke with the jsonData")
	}

	path := getPath()
	dir := filepath.Dir(path)

	_, err = os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("Unable to create dir %s, err %s", dir, err))
		}

		err = os.Mkdir(dir, os.ModePerm)

		if err != nil {
			panic(fmt.Sprintf("Unable to loac the current user, err: %s", err.Error()))
		}
	}

	err = os.WriteFile(path, jsonData, 0644)

	if err != nil {
		panic(fmt.Sprintf("Failed to Write file %s, err %s", path, err))
	}
}

func LoadData() boards.Model {
	var boards boards.Model
	path := getPath()
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
