package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func BeginWatch() {

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				// if event.Has(fsnotify.Write) {
				// 	log.Println("modified file:", event.Name)
				// }
				if event.Has(fsnotify.Create); strings.HasSuffix(event.Name, "current_run.save") {
					log.Println("Current run modified")
					// get latest map point or pre-finished room
					run, err := ReadCurrentSave()
					if err != nil {
						log.Fatal(err)
					}
					mapH := run.MapPointHistory[0]
					lastNode := mapH[len(mapH)-1]
					AnalyzeRoom(lastNode)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	log.Println("Watching on path: ", PROFILE_PATH)
	path := filepath.Join(PROFILE_PATH)
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})

}

func AnalyzeRoom(room MapPoint) {
	fmt.Println("Room:", room.MapPointType.Name())
	switch room.MapPointType {
	case MapNodeTypeAncient:
		fmt.Println("Options:")
		if ac := room.PlayerStats[0].AncientChoice; ac != nil {
			for _, c := range *ac {
				fmt.Println("Option:", c.Name())
				fmt.Println("Was chosen:", c.WasChosen)
			}
		}
	case MapNodeTypeMonster:
		fmt.Println("")
		if cc := room.PlayerStats[0].CardChoices; cc != nil {
			for _, c := range *cc {
				fmt.Println("Option:", c.Card.Name())
				fmt.Println("Was chosen:", c.WasPicked)
			}
		}
	default:
	}
}
