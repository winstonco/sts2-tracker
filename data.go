package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
)

func connDB() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	address := "127.0.0.1:3306"
	cfg.Addr = address
	cfg.DBName = "sts2_tracker"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to db at: " + address)

	return db, nil
}

// Model for cards table
type DBCardsModel struct {
	CardId      string
	CharacterId string
	FloorSeen   int
	WasPicked   bool
	RunId       int64
}

// Model for runs table
type DBRunsModel struct {
	Id                int64
	Act1              string
	Act2              string
	Act3              string
	Ascension         int
	BuildId           string
	GameMode          string
	KilledByEncounter string
	KilledByEvent     string
	RunTime           int
	SchemaVersion     int
	Seed              string
	StartTime         int
	WasAbandoned      bool
	IsWin             bool
}

// Model for card_choice_options table
type DBCardChoiceOptionsModel struct {
	Id                  int64
	CardChoiceId        int64
	CardId              string
	CurrentUpgradeLevel *int
}

// Model for map_point_card_choices table
type DBMapPointCardChoicesModel struct {
	Id     int64
	RunId  int64
	Floor  int
	Choice *string
}

// File name = start time
func getRunWithName(db *sql.DB, id int) (DBRunsModel, error) {
	var r DBRunsModel
	row := db.QueryRow("SELECT * FROM runs WHERE start_time = ?", id)
	if err := row.Scan(
		&r.Id,
		&r.Act1,
		&r.Act2,
		&r.Act3,
		&r.Ascension,
		&r.BuildId,
		&r.GameMode,
		&r.KilledByEncounter,
		&r.KilledByEvent,
		&r.RunTime,
		&r.SchemaVersion,
		&r.Seed,
		&r.StartTime,
		&r.WasAbandoned,
		&r.IsWin,
	); err != nil {
		return r, err
	}
	return r, nil
}

func insertRun(db *sql.DB, r DBRunsModel) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO runs VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.Act1, r.Act2, r.Act3, r.Ascension, r.BuildId, r.GameMode, r.KilledByEncounter, r.KilledByEvent, r.RunTime, r.SchemaVersion, r.Seed, r.StartTime, r.WasAbandoned, r.IsWin,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertCard(db *sql.DB, c DBCardsModel) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO cards VALUES (DEFAULT, ?, ?, ?, ?, ?)",
		c.CardId, c.CharacterId, c.FloorSeen, c.WasPicked, c.RunId,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertCardChoiceOption(db *sql.DB, c DBCardChoiceOptionsModel) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO card_choice_options VALUES (DEFAULT, ?, ?, ?)",
		c.CardChoiceId, c.CardId, c.CurrentUpgradeLevel,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertMapPointCardChoice(db *sql.DB, c DBMapPointCardChoicesModel) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO map_point_card_choices VALUES (DEFAULT, ?, ?, ?)",
		c.RunId, c.Floor, c.Choice,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func saveRun(db *sql.DB, run PastRunFile) error {
	// Save run in db
	rId, err := insertRun(db, DBRunsModel{
		Act1:              run.Acts[0],
		Act2:              run.Acts[1],
		Act3:              run.Acts[2],
		Ascension:         run.Ascension,
		BuildId:           run.BuildId,
		GameMode:          run.GameMode,
		KilledByEncounter: run.KilledByEncounter,
		KilledByEvent:     run.KilledByEvent,
		RunTime:           run.RunTime,
		SchemaVersion:     run.SchemaVersion,
		Seed:              run.Seed,
		StartTime:         run.StartTime,
		WasAbandoned:      run.WasAbandoned,
		IsWin:             run.Win,
	})
	if err != nil {
		return err
	}
	log.Printf("Added new run: %d\n", rId)
	// make map of player id to character for later use
	pIdToChar := make(map[int]string)
	for _, p := range run.Players {
		pIdToChar[p.Id] = p.Character
	}
	// Read map_point_history
	// For each act
	toInsert := make([]DBCardChoiceOptionsModel, 0)
	floor := 0
	for _, actPoints := range run.MapPointHistory {
		// For each map point
		for _, mapPoint := range actPoints {
			floor++
			// For each player
			for _, player := range mapPoint.PlayerStats {
				// For each card choice
				if player.CardChoices != nil {
					// put each choice in a list to add after adding map_point_card_choices row, which requires knowing the card choice made (or skip)
					var choice *string = nil
					for _, choice_opt := range *player.CardChoices {
						card := choice_opt.Card
						if choice_opt.WasPicked {
							choice = &card.Id
						}
						choiceOpt := DBCardChoiceOptionsModel{
							CardId:              card.Id,
							CurrentUpgradeLevel: card.CurrentUpgradeLevel,
						}
						log.Println(choiceOpt)
						toInsert = append(toInsert, choiceOpt)
						// // Save card data in db
						// cId, err := insertCard(db, DBCardsModel{
						// 	CardId:      card.Id,
						// 	CharacterId: pIdToChar[player.PlayerId],
						// 	FloorSeen:   floor,
						// 	WasPicked:   choice.WasPicked,
						// 	RunId:       rId,
						// })
						// if err != nil {
						// 	return err
						// }
					}
					// insert map point card choice parent row
					mpccId, err := insertMapPointCardChoice(db, DBMapPointCardChoicesModel{
						RunId:  rId,
						Floor:  floor,
						Choice: choice,
					})
					if err != nil {
						return err
					}
					log.Printf("Added new map_point_card_choices: %d\n", mpccId)
					// insert all card choice child rows
					for len(toInsert) > 0 {
						cc := toInsert[0]
						_, err := insertCardChoiceOption(db, DBCardChoiceOptionsModel{
							CardChoiceId:        mpccId,
							CardId:              cc.CardId,
							CurrentUpgradeLevel: cc.CurrentUpgradeLevel,
						})
						if err != nil {
							return err
						}
						log.Printf("Added new card: %s\n", cc.CardId)
						toInsert = toInsert[1:]
					}
				}
			}
		}
	}
	return nil
}

func readAndSaveRunHistory(db *sql.DB) {
	log.Println("Reading local run history")
	// Get runs in game history
	dirPath := filepath.Join(PROFILE_SAVES_PATH, "history")
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, dirEntry := range files {
		log.Println(dirEntry.Name())
		fp := filepath.Join(dirPath, dirEntry.Name())
		file, err := os.ReadFile(fp)
		if err != nil {
			log.Fatal(err)
		}
		run, err := readPastRunFile(file)
		// Check if run exists
		_, err = getRunWithName(db, run.StartTime)
		switch err {
		case nil:
			// already saved
			log.Println("Run exists")
		case sql.ErrNoRows:
			// add run
			if err := saveRun(db, run); err != nil {
				log.Fatalf("\n%v", err)
			}
		default:
			log.Fatal(err)
		}
	}
}
