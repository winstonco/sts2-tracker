package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

func ReadCurrentSave() (CurrentRunSave, error) {
	path := filepath.Join(STEAM_USERDATA_PATH, STEAM_USER_ID, GAME_ID, "remote", PROFILE, "saves", "current_run.save")

	file, err := os.ReadFile(path)
	if err != nil {
		return CurrentRunSave{}, err
	}

	currentRunSave := CurrentRunSave{}
	err = json.Unmarshal(file, &currentRunSave)
	if err != nil {
		return CurrentRunSave{}, err
	}
	return currentRunSave, nil
}
