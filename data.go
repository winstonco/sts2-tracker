package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

func ReadCurrentSave() (CurrentRunSaveFile, error) {
	file, err := os.ReadFile(CURRENT_RUN_SAVE_PATH)
	if err != nil {
		return CurrentRunSaveFile{}, err
	}

	currentRunSave := CurrentRunSaveFile{}
	if err := json.Unmarshal(file, &currentRunSave); err != nil {
		return CurrentRunSaveFile{}, err
	}
	return currentRunSave, nil
}

func ReadProgressSave() (ProgressSaveFile, error) {
	file, err := os.ReadFile(PROGRESS_SAVE_PATH)
	if err != nil {
		return ProgressSaveFile{}, err
	}

	progressSave := ProgressSaveFile{}
	if err := json.Unmarshal(file, &progressSave); err != nil {
		return ProgressSaveFile{}, err
	}
	return progressSave, nil
}

func ReadPastRun(path string) (PastRunFile, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return PastRunFile{}, err
	}

	pastRun := PastRunFile{}
	if err := json.Unmarshal(file, &pastRun); err != nil {
		return PastRunFile{}, err
	}
	return pastRun, nil
}

func ReadAndSavePastRun(file *excelize.File, runId string) error {
	path := filepath.Join(PROFILE_SAVES_PATH, "history", runId)
	run, err := ReadPastRun(path)
	if err != nil {
		return err
	}
	// Get the num of rows in cards table
	b2, err := file.GetCellValue("Metadata", "B2")
	if err != nil {
		return err
	}
	cardsSaved, err := strconv.Atoi(b2)
	if err != nil {
		return err
	}
	// make map of player id to character for later use
	pIdToChar := make(map[int]string)
	for _, p := range run.Players {
		pIdToChar[p.Id] = p.Character
	}
	// Read map_point_history
	// For each act
	for i, actPoints := range run.MapPointHistory {
		// For each map point
		for j, mapPoint := range actPoints {
			// calc floor number (floor in act + floor counts of prev acts)
			floor := j + 1
			for k := 0; k < i-1; k++ {
				floor += len(run.MapPointHistory[i])
			}
			// For each player
			for _, player := range mapPoint.PlayerStats {
				// For each card choice
				if player.CardChoices != nil {
					for _, choice := range *player.CardChoices {
						card := choice.Card
						file.SetSheetRow("Cards", fmt.Sprintf("A%d", cardsSaved+2), &[]any{
							card.Id,
							pIdToChar[player.PlayerId],
							run.Ascension,
							floor,
							choice.WasPicked,
							run.Acts[0],
							run.Acts[1],
							run.Acts[2],
							run.Win,
							runId,
							run.BuildId,
						})
						cardsSaved++
						file.SetCellValue("Metadata", "B2", cardsSaved)
					}
				}
			}
		}
	}
	if err := file.Save(); err != nil {
		return err
	}
	return nil
}

const SPIRE_DATA_VERSION = "1"

func makeNewSpireDataFile() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := f.SetSheetName("Sheet1", "Metadata"); err != nil {
		return err
	}
	if err := f.SetSheetRow("Metadata", "A1", &[]any{
		"SpireDataVersion",
		SPIRE_DATA_VERSION,
	}); err != nil {
		return err
	}
	if err := f.SetSheetRow("Metadata", "A2", &[]any{
		"CardsSaved",
		0,
	}); err != nil {
		return err
	}
	if _, err := f.NewSheet("Runs"); err != nil {
		return err
	}
	if _, err := f.NewSheet("Cards"); err != nil {
		return err
	}
	if err := f.SetSheetRow("Cards", "A1", &[]any{
		"Card ID",
		"Character ID",
		"Ascension",
		"Floor Seen",
		"Picked",
		"Act1",
		"Act2",
		"Act3",
		"Won",
		"Run ID",
		"Version",
	}); err != nil {
		return err
	}
	if err := f.SaveAs("SpireData.xlsx"); err != nil {
		return err
	}
	return nil
}

func AggregateRunData() {
	if !FileExists("./SpireData.xlsx") {
		if err := makeNewSpireDataFile(); err != nil {
			log.Fatal(err)
		}
	}
	f, err := excelize.OpenFile("SpireData.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	spireDataVersion, err := f.GetCellValue("Metadata", "B1")
	if err != nil {
		log.Fatal(err)
	}
	// If SpireData is outdated, clear history of runs already collected
	if spireDataVersion != SPIRE_DATA_VERSION {
		if err := f.RemoveCol("Runs", "A"); err != nil {
			log.Fatal(err)
		}
		if err := f.Save(); err != nil {
			log.Fatal(err)
		}
	}
	// Get collected collRuns list
	var collRuns []string
	cols, err := f.Cols("Runs")
	if err != nil {
		log.Fatal(err)
	}
	if cols.Next() {
		collRuns, err = cols.Rows()
		if err != nil {
			log.Fatal(err)
		}
	}
	// Get runs in game history
	files, err := os.ReadDir(filepath.Join(PROFILE_SAVES_PATH, "history"))
	if err != nil {
		log.Fatal(err)
	}
	// Compare uncollected runs
	// TODO: Maybe this is too slow
	countNew := 0
	for _, file := range files {
		// Check if in collected runs list
		collected := false
		for _, r := range collRuns {
			if r == file.Name() {
				collected = true
			}
		}
		if !collected {
			countNew++
			if err := ReadAndSavePastRun(f, file.Name()); err != nil {
				log.Fatal(err)
			}
			f.SetCellValue("Runs", fmt.Sprintf("A%d", len(collRuns)+countNew), file.Name())
			if err := f.Save(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

/*
func WriteExcelExample() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal()
		}
	}()
	// Create a new sheet.
	// index, err := f.NewSheet("Sheet1")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// f.SetActiveSheet(index)
	if err := f.SaveAs("SpireData.xlsx"); err != nil {
		log.Fatal(err)
	}
}

func ReadExcelExample() {
	f, err := excelize.OpenFile("SpireData.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		for _, colCell := range row {
			log.Print(colCell, "\t")
		}
		log.Println()
	}
}
*/
