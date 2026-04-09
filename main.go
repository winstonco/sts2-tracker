package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const STEAM_USERDATA_PATH = `C:\Program Files (x86)\Steam\userdata\`
const STEAM_USER_ID = "286123118"
const GAME_ID = "2868840"
const PROFILE = "profile1"

func main() {
	log.Println("Hello World")

	// // Create new watcher.
	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer watcher.Close()

	// // Start listening for events.
	// go func() {
	// 	for {
	// 		select {
	// 		case event, ok := <-watcher.Events:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println("event:", event)
	// 			if event.Has(fsnotify.Write) {
	// 				log.Println("modified file:", event.Name)
	// 			}
	// 		case err, ok := <-watcher.Errors:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println("error:", err)
	// 		}
	// 	}
	// }()

	// // Add a path.
	// err = watcher.Add(SPIRE_PATH)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// <-make(chan struct{})

	// Read progress.save
	path := filepath.Join(STEAM_USERDATA_PATH, STEAM_USER_ID, GAME_ID, "remote", PROFILE, "saves", "progress.save")
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	progressSave := ProgressSave{}

	if err := json.Unmarshal(file, &progressSave); err != nil {
		log.Fatal(err)
	}

	for _, c := range progressSave.CharacterStats {
		fmt.Println(c.Name())
		fmt.Printf("Winrate: %.2f\n", float64(c.TotalWins)/(float64(c.TotalWins)+float64(c.TotalLosses)))
	}

	// Read current_run.save
	path = filepath.Join(STEAM_USERDATA_PATH, STEAM_USER_ID, GAME_ID, "remote", PROFILE, "saves", "current_run.save")
	file, err = os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	currentRunSave := CurrentRunSave{}

	if err := json.Unmarshal(file, &currentRunSave); err != nil {
		log.Fatal(err)
	}
	fmt.Println(currentRunSave)
}
