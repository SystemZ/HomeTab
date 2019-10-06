package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	var gamesAdded int
	for _, game := range export {
		log.Printf("%v", game.Name)
		gamesAdded++
	}
	log.Printf("Games added: %v", gamesAdded)
}
