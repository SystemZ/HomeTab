package kernelcare

import (
	"encoding/json"
	"github.com/gorilla/feeds"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type KernelcareList map[string]Kernel

type Kernel struct {
	Name        string  `json:"name"`
	Timestamp   float64 `json:"timestamp"`
	KernelUname string  `json:"kernel_uname"`
	Time        string  `json:"time"`
	Flavor      string  `json:"flavor"`
	Latest      int     `json:"latest"`
}

func Fetch(flavorFilter string) *feeds.Feed {

	feed := &feeds.Feed{
		Title:       "Kernelcare",
		Link:        &feeds.Link{Href: "https://patches.kernelcare.com/"},
		Description: "Latest kernel patches",
		//Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		//Created:     now,
	}

	url := "https://patches.kernelcare.com/patches.json"
	spaceClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	//req.Header.Set("User-Agent", "spacecount-tutorial")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	people1 := KernelcareList{}
	jsonErr := json.Unmarshal(body, &people1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for kernelHash, kernel := range people1 {
		if flavorFilter != kernel.Flavor {
			continue
		}
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          kernel.Name,
			Title:       kernel.Name,
			Link:        &feeds.Link{Href: "https://patches.kernelcare.com/" + kernelHash + "/" + strconv.Itoa(kernel.Latest) + "/kpatch.html"},
			Description: kernel.KernelUname,
			//Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created: time.Unix(int64(kernel.Timestamp), 0),
			Updated: time.Unix(int64(kernel.Timestamp), 0),
		})

	}
	return feed
}

func Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	feed := Fetch("pve-6")
	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.WriteString(w, atom)
	if err != nil {
		log.Fatal(err)
	}
}
