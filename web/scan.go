package web

import (
	"encoding/json"
	"gitlab.com/systemz/gotag/core"
	"gitlab.com/systemz/gotag/model"
	"log"
	"net/http"
)

type ScanReq struct {
	Path string `json:"path"`
}

func Scan(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get tag name from JSON body
	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var scan []ScanReq
	decoder.Decode(&scan)

	if len(scan) < 1 {
		log.Printf("No paths to scan!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		for _, toScan := range scan {
			// ignore empty
			if len(toScan.Path) < 1 {
				continue
			}
			core.ScanMono(model.DB, toScan.Path, int(userInfo.Id))
		}
	}()
}
