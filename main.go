package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const STEAM_USERDATA_PATH = `C:\Program Files (x86)\Steam\userdata\`
const STEAM_USER_ID = "286123118"
const GAME_ID = "2868840"
const PROFILE = "profile1"

var PROFILE_PATH = filepath.Join(STEAM_USERDATA_PATH, STEAM_USER_ID, GAME_ID, "remote", PROFILE, "saves")

func main() {
	log.Println("Hello World")

	// Read progress.save
	path := filepath.Join(PROFILE_PATH, "progress.save")
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

	if FileExists(path) {
		run, err := ReadCurrentSave()
		if err != nil {
			log.Fatal(err)
		}
		player := run.Players[0]
		fmt.Println("Current run:")
		fmt.Print(player.Name())
		fmt.Println(" A" + strconv.Itoa(run.Ascension))
		fmt.Println("Deck:")
		for _, c := range player.Deck {
			fmt.Println(c.Name())
		}
		fmt.Println("Relics:")
		for _, r := range player.Relics {
			fmt.Println(r.Name())
		}
	}

	// Start watching
	BeginWatch()
}
