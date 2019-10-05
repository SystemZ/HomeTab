package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"
)

// store game sessions
// this log doesn't have entries pre 2012-10-03 -
// instead it's counted in time_p and time_s as a summary
type ZfireLog []struct {
	Name  string    `json:"name"`
	TimeP int       `json:"time_p"`
	TimeS int       `json:"time_s"`
	Date  time.Time `json:"date"`
}

// sorting functions
func (a ZfireLog) Len() int           { return len(a) }
func (a ZfireLog) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ZfireLog) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// store complete game list with game time summary
type ZfireGame []struct {
	Date  time.Time `json:"date"`
	Name  string    `json:"name"`
	TimeP int       `json:"time_p"`
	TimeS int       `json:"time_s"`
	Tags  []string  `json:"tags"`
}

// sorting functions
func (a ZfireGame) Len() int           { return len(a) }
func (a ZfireGame) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ZfireGame) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// temp local store
type ZfireLogGameList struct {
	Name string
	Date time.Time
}

// helper
func IsInList(list []ZfireLogGameList, name string) bool {
	for _, v := range list {
		if v.Name == name {
			return true
		}
	}
	return false
}

// helper
func IsInArray(list []string, name string) bool {
	for _, v := range list {
		if v == name {
			return true
		}
	}
	return false
}

func ImportZfire(pathToJson string) {
	// 1. Load
	//

	//read zfire game log from disk
	zfireLogRaw, err := ioutil.ReadFile(pathToJson + "/log.json")
	if err != nil {
		log.Printf("%v", err)
	}
	//parse zfire game log from json
	var zfireLog ZfireLog
	err = json.Unmarshal(zfireLogRaw, &zfireLog)
	if err != nil {
		fmt.Println("error:", err)
	}

	//read zfire game list from disk
	zfireGamesRaw, err := ioutil.ReadFile(pathToJson + "/games.json")
	if err != nil {
		log.Printf("%v", err)
	}
	//parse zfire game log from json
	var zfireGames ZfireGame
	err = json.Unmarshal(zfireGamesRaw, &zfireGames)
	if err != nil {
		fmt.Println("error:", err)
	}

	// 2. Sort
	//

	// sort zfire game log from oldest to newest entries
	sort.Sort(zfireLog)
	sort.Sort(zfireGames)

	// 3.Analyze
	//

	// check for game duplicates
	var gameListDuplicates []string
	var zfireGamesParsedTmp []string
	// detect duplicate/cross platform games
	log.Printf("%v", "Checking for duplicates...")
	for _, game := range zfireGames {
		if IsInArray(zfireGamesParsedTmp, game.Name) {
			gameListDuplicates = append(gameListDuplicates, game.Name)
			log.Printf("dup: %v", game.Name)
		} else {
			zfireGamesParsedTmp = append(zfireGamesParsedTmp, game.Name)
		}
	}

	//check game list for lacking creation date
	log.Printf("%v", "Checking dates in game list...")
	var gameListLackingCreationDates []string
	for _, game := range zfireGames {
		if game.Date.IsZero() {
			gameListLackingCreationDates = append(gameListLackingCreationDates, game.Name)
		}
	}

	// various magic operations
	//
	var uniqueGameList []ZfireLogGameList
	// get earliest session of game
	for _, game := range zfireLog {

		if game.Name == "Zeno Clash" {
			log.Printf("%v", "-")
		}

		if IsInList(uniqueGameList, game.Name) {
			continue
		}

		var earliestDate time.Time
		//log.Printf("%v", game.Name)
		earliestDate = game.Date
		for _, gamez := range zfireLog {
			testedDate := gamez.Date
			if gamez.Date.Before(testedDate) {
				earliestDate = gamez.Date
			}
		}
		//log.Printf("%v", earliestDate)
		uniqueGameList = append(uniqueGameList, ZfireLogGameList{Name: game.Name, Date: earliestDate})

	}

	var differenceGamesList []string
	var isOnList bool
	for _, game := range zfireGames {
		isOnList = false
		for _, gamez := range uniqueGameList {
			if game.Name == gamez.Name {
				isOnList = true
			}
		}
		if !isOnList {
			differenceGamesList = append(differenceGamesList, game.Name)
		}
	}

	log.Printf("Games duplicated: %v", len(gameListDuplicates))
	//log.Printf("Games with log: %v", len(uniqueGameList))
	log.Printf("Games without log: %v", len(differenceGamesList))
	log.Printf("All games: %v", len(zfireGames))
	//log.Printf("%v", differenceGamesList)
	//log.Printf("%v", gameListDuplicates)

	// all games in Zfire: 678

	// export converted json
	//

	zfireGameListByte, err := json.MarshalIndent(zfireGames, "", "  ")
	ioutil.WriteFile("export/zfireGames.json", zfireGameListByte, 0644)
	zfireLogByte, err := json.MarshalIndent(zfireLog, "", "  ")
	ioutil.WriteFile("export/zfireLog.json", zfireLogByte, 0644)

	diffGameListByte, err := json.MarshalIndent(differenceGamesList, "", "  ")
	ioutil.WriteFile("export/zfireGamesWithoutLogs.json", diffGameListByte, 0644)
	gameListLackingCreationDatesByte, err := json.MarshalIndent(gameListLackingCreationDates, "", "  ")
	ioutil.WriteFile("export/zfireGamesLackingCreationDates.json", gameListLackingCreationDatesByte, 0644)

	//var counter Counter
	//counter.Name = game.Name
	////FIXME timezones
	//counter.CreatedAt = &game.Date
	//counter.UpdatedAt = &game.Date
	//
	//err := DB.Save(&counter).Error
	//if err != nil {
	//	log.Printf("%v", err)
	//}
}
