package tm_poznan

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gorilla/feeds"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Fetch() *feeds.Feed {

	feed := &feeds.Feed{
		Title:       "TM Poznan",
		Link:        &feeds.Link{Href: "http://teatr-muzyczny.poznan.pl/"},
		Description: "Latest plays",
	}

	url := "http://teatr-muzyczny.poznan.pl/api/shows/calendar?page=1&per_page=100&from_now=1"
	// this site have some problems with TLS, let's ignore TLS problems for now
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	spaceClient := http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
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

	people1 := Welcome{}
	jsonErr := json.Unmarshal(body, &people1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for _, play := range people1.Data {
		for _, playInfo := range play.Data {
			for _, playDates := range playInfo.Relationships.Callendar.Data {

				// parsing date
				layout := "2006-01-02 15:04:05 -0700"
				dateString := playDates.Attributes.Publication + " +0200"
				t, err := time.Parse(layout, dateString)
				if err != nil {
					log.Println("Error while parsing date :", err)
				}

				// add item to feed
				feed.Items = append(feed.Items, &feeds.Item{
					Id:          playInfo.Attributes.Title + playDates.Attributes.Publication,
					Title:       playInfo.Attributes.Title + " " + playDates.Attributes.Publication,
					Link:        &feeds.Link{Href: "http://teatr-muzyczny.poznan.pl/"},
					Description: playInfo.Attributes.Title + "<br><img src='" + playInfo.Attributes.ImgURL + "'/>",
					Created:     t,
					Updated:     t,
				})
			}
		}
	}
	return feed
}

func Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	feed := Fetch()
	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.WriteString(w, atom)
	if err != nil {
		log.Fatal(err)
	}
}
