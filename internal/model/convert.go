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

// final export json
type TaskTabExport struct {
	Name            string               `json:"name"`
	CreatedAt       time.Time            `json:"createdAt"`
	Tags            []string             `json:"tags"`
	ZfireTimeSumS   uint                 `json:"zfireTimeSummaryS"`
	TaskTabTimeSumS uint                 `json:"taskTabTimeSummaryS"`
	ZfireTimeSumP   uint                 `json:"zfireTimeSummaryP"`
	TaskTabTimeSumP uint                 `json:"taskTabTimeSummaryP"`
	SessionsS       TaskTabExportSession `json:"sessionsS"`
	SessionsP       TaskTabExportSession `json:"sessionsP"`
}
type TaskTabExports []TaskTabExport

func (a TaskTabExports) Len() int           { return len(a) }
func (a TaskTabExports) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TaskTabExports) Less(i, j int) bool { return a[i].CreatedAt.Before(a[j].CreatedAt) }

type TaskTabExportSession []struct {
	Duration  int
	StartedAt time.Time
	EndedAt   time.Time
	Precise   uint
}

func (a TaskTabExportSession) Len() int           { return len(a) }
func (a TaskTabExportSession) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TaskTabExportSession) Less(i, j int) bool { return a[i].StartedAt.Before(a[j].StartedAt) }

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
		if IsInList(uniqueGameList, game.Name) {
			continue
		}

		var earliestDate time.Time
		earliestDate = game.Date
		for _, gamez := range zfireLog {
			testedDate := gamez.Date
			if gamez.Date.Before(testedDate) {
				earliestDate = gamez.Date
			}
		}
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

	var finalExport []TaskTabExport
	var finalExportErrors int
	var finalExportRecounstructed int
	// Export
	for _, game := range zfireGames {
		var skipGame bool

		// skip empty games
		gamesToSkip := []string{
			"Warhammer 40,000: Dawn of War GOTY",
			"Toki Tori 2+",
			"Halo 4",
		}
		for _, emptyGame := range gamesToSkip {
			if game.Name == emptyGame {
				skipGame = true
			}
		}

		// skip duplicated games for now
		// we resolve platform tags later
		for _, dupGame := range gameListDuplicates {
			if game.Name == dupGame {
				skipGame = true
				log.Printf("dup: %v", game.Name)
			}
		}

		// skip game for various reasons
		if skipGame {
			continue
		}

		var sessionsS TaskTabExportSession
		var sessionsP TaskTabExportSession
		var taskTabTimeSumS uint
		var taskTabTimeSumP uint
		for _, gameInLog := range zfireLog {
			if gameInLog.Name != game.Name {
				continue
			}

			// work on proper and precise game logs
			if gameInLog.TimeS != 0 {
				startedAt := gameInLog.Date.Add(-time.Second * time.Duration(gameInLog.TimeS))
				duration := gameInLog.Date.Sub(startedAt).Seconds()
				sessionsS = append(sessionsS, struct {
					Duration  int
					StartedAt time.Time
					EndedAt   time.Time
					Precise   uint
				}{Duration: int(duration), StartedAt: startedAt, EndedAt: gameInLog.Date, Precise: 1})
				taskTabTimeSumS += uint(gameInLog.Date.Sub(startedAt).Seconds())
			}
			if gameInLog.TimeP != 0 {
				startedAt := gameInLog.Date.Add(-time.Second * time.Duration(gameInLog.TimeP))
				duration := gameInLog.Date.Sub(startedAt).Seconds()
				sessionsP = append(sessionsP, struct {
					Duration  int
					StartedAt time.Time
					EndedAt   time.Time
					Precise   uint
				}{Duration: int(duration), StartedAt: startedAt, EndedAt: gameInLog.Date, Precise: 1})
				taskTabTimeSumP += uint(gameInLog.Date.Sub(startedAt).Seconds())
			}
		}

		// reconstruction of time
		//add one big session ending at 01.10.2012
		//first real logs starts at 03.10.2012 so this prevent conflicts
		warsaw, _ := time.LoadLocation("Europe/Warsaw")
		zfireTimestamp := time.Date(2012, 10, 1, 0, 0, 0, 0, warsaw)
		if taskTabTimeSumS != uint(game.TimeS) {
			startedAt := zfireTimestamp.Add(-time.Second * time.Duration(game.TimeS-int(taskTabTimeSumS)))
			duration := zfireTimestamp.Sub(startedAt).Seconds()
			sessionsS = append(sessionsS, struct {
				Duration  int
				StartedAt time.Time
				EndedAt   time.Time
				Precise   uint
			}{Duration: int(duration), StartedAt: startedAt, EndedAt: zfireTimestamp, Precise: 0})
			taskTabTimeSumS += uint(zfireTimestamp.Sub(startedAt).Seconds())
			finalExportRecounstructed++
		}
		if taskTabTimeSumP != uint(game.TimeP) {
			startedAt := zfireTimestamp.Add(-time.Second * time.Duration(game.TimeP-int(taskTabTimeSumP)))
			duration := zfireTimestamp.Sub(startedAt).Seconds()
			sessionsP = append(sessionsP, struct {
				Duration  int
				StartedAt time.Time
				EndedAt   time.Time
				Precise   uint
			}{Duration: int(duration), StartedAt: startedAt, EndedAt: zfireTimestamp, Precise: 0})
			taskTabTimeSumP += uint(zfireTimestamp.Sub(startedAt).Seconds())
			finalExportRecounstructed++
		}

		// sort sessions after reconstruction
		sort.Sort(sessionsP)
		sort.Sort(sessionsS)

		// reconstruct game creation date from oldest log entry
		if game.Date.IsZero() && (sessionsS.Len() > 0 || sessionsP.Len() > 0) {
			if sessionsS.Len() > sessionsP.Len() {
				game.Date = sessionsS[0].StartedAt
			} else if sessionsP.Len() > sessionsS.Len() {
				game.Date = sessionsP[0].StartedAt
			}
		}
		// for games never played on zfire and created on firez, set constant date as a fallback
		if game.Date.IsZero() {
			game.Date = zfireTimestamp
		}

		// write all game sessions
		dataToWrite := TaskTabExport{
			Name:            game.Name,
			CreatedAt:       game.Date,
			Tags:            game.Tags,
			ZfireTimeSumS:   uint(game.TimeS),
			ZfireTimeSumP:   uint(game.TimeP),
			TaskTabTimeSumS: taskTabTimeSumS,
			TaskTabTimeSumP: taskTabTimeSumP,
			SessionsS:       sessionsS,
			SessionsP:       sessionsP,
		}
		finalExport = append(finalExport, dataToWrite)

		// warning
		if dataToWrite.CreatedAt.IsZero() || taskTabTimeSumP != uint(game.TimeP) || taskTabTimeSumS != uint(game.TimeS) {
			log.Printf("Errors in %v", game.Name)
			finalExportErrors++
		}
	}

	//final sorting by game creation date
	sort.Sort(TaskTabExports(finalExport))

	// show stats
	log.Printf("= Stats =")
	log.Printf("Errors with time: %v, Duplicates: %v, No logs: %v", finalExportErrors, len(gameListDuplicates), len(differenceGamesList))
	log.Printf("Games exported without errors: %v/%v", len(finalExport)-finalExportErrors, len(zfireGames))
	var percentage float32
	percentage = (float32(len(finalExport)) - float32(finalExportErrors)) / float32(len(zfireGames))
	log.Printf("Integrity %v%%", percentage*100)

	// export converted json
	zfireGameListByte, err := json.MarshalIndent(zfireGames, "", "  ")
	ioutil.WriteFile("export/zfireGames.json", zfireGameListByte, 0644)

	// export json helpers
	zfireLogByte, err := json.MarshalIndent(zfireLog, "", "  ")
	ioutil.WriteFile("export/zfireLog.json", zfireLogByte, 0644)
	diffGameListByte, err := json.MarshalIndent(differenceGamesList, "", "  ")
	ioutil.WriteFile("export/zfireGamesWithoutLogs.json", diffGameListByte, 0644)
	gameListLackingCreationDatesByte, err := json.MarshalIndent(gameListLackingCreationDates, "", "  ")
	ioutil.WriteFile("export/zfireGamesLackingCreationDates.json", gameListLackingCreationDatesByte, 0644)
	finalExportByte, err := json.MarshalIndent(finalExport, "", "  ")
	ioutil.WriteFile("export/export.json", finalExportByte, 0644)
}
