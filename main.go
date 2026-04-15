package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
)

const STEAM_USERDATA_PATH = `C:\Program Files (x86)\Steam\userdata\`
const STEAM_USER_ID = "286123118"
const GAME_ID = "2868840"
const PROFILE = "profile1"

var PROFILE_SAVES_PATH = filepath.Join(STEAM_USERDATA_PATH, STEAM_USER_ID, GAME_ID, "remote", PROFILE, "saves")
var PROGRESS_SAVE_PATH = filepath.Join(PROFILE_SAVES_PATH, "progress.save")
var CURRENT_RUN_SAVE_PATH = filepath.Join(PROFILE_SAVES_PATH, "current_run.save")

func main() {
	// Read progress.save
	progressSave, err := ReadProgressSave()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range progressSave.CharacterStats {
		fmt.Println(c.Name())
		fmt.Printf("Winrate: %.2f\n", float64(c.TotalWins)/(float64(c.TotalWins)+float64(c.TotalLosses)))
	}

	// Read current_run.save
	if FileExists(CURRENT_RUN_SAVE_PATH) {
		log.Println("Current run found")
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
	} else {
		log.Println("No current run")
	}

	// Aggregate past run data
	AggregateRunData()

	// Start watching
	// BeginWatch()
}
