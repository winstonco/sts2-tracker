package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	db, err := connDB()
	if err != nil {
		log.Fatal(err)
	}

	// Read progress.save
	progressSave, err := readProgressSave()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range progressSave.CharacterStats {
		fmt.Println(c.Name())
		fmt.Printf("Winrate: %.2f\n", float64(c.TotalWins)/(float64(c.TotalWins)+float64(c.TotalLosses)))
	}

	// Read current_run.save
	_, err = os.Stat(CURRENT_RUN_SAVE_PATH)
	if err == nil { // file exists
		log.Println("Current run found")
		run, err := readCurrentSave()
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

	// Save past run data
	readAndSaveRunHistory(db)

	// Start watching
	// BeginWatch()
}
