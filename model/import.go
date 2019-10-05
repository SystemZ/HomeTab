package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type ZfireGame struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Name  string   `json:"name"`
	TimeP int      `json:"time_p"`
	TimeS int      `json:"time_s"`
	Tags  []string `json:"tags"`
}

type ZfireLog []struct {
	Name  string    `json:"name"`
	TimeP int       `json:"time_p"`
	TimeS int       `json:"time_s"`
	Date  time.Time `json:"date"`
}

type ZfireLogGameList struct {
	Name string
	Date time.Time
}

func IsInList(list []ZfireLogGameList, name string) bool {
	for _, v := range list {
		if v.Name == name {
			return true
		}
	}
	return false
}

func IsInArray(list []string, name string) bool {
	for _, v := range list {
		if v == name {
			return true
		}
	}
	return false
}

func ImportZfire(pathToJson string) {
	// load data //
	//
	// parse zfireLog
	zfireLogRaw, err := ioutil.ReadFile(pathToJson + "/log.json")
	if err != nil {
		log.Printf("%v", err)
	}
	var zfireLog ZfireLog
	err = json.Unmarshal(zfireLogRaw, &zfireLog)
	if err != nil {
		fmt.Println("error:", err)
	}
	// parse zfireGames
	zfireGamesRaw, err := ioutil.ReadFile(pathToJson + "/games.json")
	if err != nil {
		log.Printf("%v", err)
	}
	var zfireGames []ZfireGame
	err = json.Unmarshal(zfireGamesRaw, &zfireGames)
	if err != nil {
		fmt.Println("error:", err)
	}

	var zfireGamesDuplicate []string
	var zfireGamesParsed []string
	// detect duplicate/cross platform games
	log.Printf("%v", "Checking for duplicates...")
	for _, game := range zfireGames {
		if IsInArray(zfireGamesParsed, game.Name) {
			zfireGamesDuplicate = append(zfireGamesDuplicate, game.Name)
			log.Printf("dup: %v", game.Name)
		} else {
			zfireGamesParsed = append(zfireGamesParsed, game.Name)
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

	log.Printf("Games duplicated: %v", len(zfireGamesDuplicate))
	log.Printf("Games with log: %v", len(uniqueGameList))
	log.Printf("Games without log: %v", len(differenceGamesList))
	log.Printf("All games: %v", len(zfireGames))
	//log.Printf("%v", differenceGamesList)
	//log.Printf("%v", zfireGamesDuplicate)

	// all games in Zfire: 678

	diffGameListByte, err := json.MarshalIndent(differenceGamesList, "", "  ")
	ioutil.WriteFile("export/gameList.json", diffGameListByte, 0644)
}
