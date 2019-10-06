package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func ImportCountersFromJson(pathToFile string) {
	//readfrom disk
	exportRaw, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		log.Printf("%v", err)
	}
	//parse from json
	var export TaskTabExports
	err = json.Unmarshal(exportRaw, &export)
	if err != nil {
		fmt.Println("error:", err)
	}

	// add counter (game)
	var gamesAdded int
	var gameSessionAdded int
	for _, game := range export {
		//log.Printf("Adding %v ...", game.Name)
		now := time.Now()
		newGame := Counter{
			Name:      game.Name,
			ProjectId: 1,
			CreatedAt: &game.CreatedAt,
			UpdatedAt: &now,
		}
		res := DB.Save(&newGame)
		if res.Error != nil {
			log.Printf("Error adding %v", game.Name)
		}
		// add counter session (game session) for P
		for _, gameSession := range game.SessionsP {
			newSession := CounterSession{
				CounterId: newGame.Id,
				UserId:    2,
				StartedAt: &gameSession.StartedAt,
				EndedAt:   &gameSession.EndedAt,
				CreatedAt: &gameSession.StartedAt,
				UpdatedAt: &now,
			}
			res := DB.Save(&newSession)
			if res.Error != nil {
				log.Printf("Error adding session to %v", game.Name)
			} else {
				gameSessionAdded++
			}
		}
		// add counter session (game session) for S
		for _, gameSession := range game.SessionsS {
			newSession := CounterSession{
				CounterId: newGame.Id,
				UserId:    1,
				StartedAt: &gameSession.StartedAt,
				EndedAt:   &gameSession.EndedAt,
				CreatedAt: &gameSession.StartedAt,
				UpdatedAt: &now,
			}
			res := DB.Save(&newSession)
			if res.Error != nil {
				log.Printf("Error adding session to %v", game.Name)
			} else {
				gameSessionAdded++
			}
		}

		gamesAdded++
	}
	log.Printf("Games added: %v", gamesAdded)
	log.Printf("Sessions added: %v", gameSessionAdded)
}
