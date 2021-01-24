package helios

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HeliosResponse struct {
	Html string `json:"html"`
}

type HeliosMovie struct {
	Day         int
	Title       string
	Description string
	Url         string
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

func Fetch() *feeds.Feed {
	feed := &feeds.Feed{
		Title:       "Helios CHR Kupiec",
		Link:        &feeds.Link{Href: "https://www.helios.pl/23,Szczecin/Repertuar/"},
		Description: "Latest movies",
	}
	var movies []HeliosMovie

	// get info from all 14 days including today, set less for debugging
	for day := 0; day <= 13; day++ {
		doc := GetMoviesFromDay(day)
		doc.Find(".seance").Each(func(i int, s *goquery.Selection) {

			//get basic info
			title := strings.TrimSpace(s.Find(".movie-title").Text())
			description := strings.TrimSpace(s.Find(".movie-description").Text())
			url, _ := s.Find(".movie-link").Attr("href")

			// get all hours available
			var hours []string
			s.Find(".hour-link").Each(func(i int, hour *goquery.Selection) {
				hours = append(hours, hour.Text())
			})

			duplicate := false
			// prevent adding same title of movie to feed again
			for _, v := range movies {
				if v.Title == title {
					duplicate = true
					break
				}
			}

			// this is not a duplicate so append to feed freely
			if !duplicate {
				//movieDate := timeNow.AddDate(0, 0, movie.Day).Format("2006-01-02")
				//FIXME parse date from site and make it immutable
				// some time magic

				timeNow := time.Now().Local()
				midnight := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.Local)

				currentYear := midnight.Format("2006")
				movieDate := midnight.AddDate(0, 0, day)
				movies = append(movies, HeliosMovie{Title: title})
				feed.Items = append(feed.Items, &feeds.Item{
					Id:          currentYear + "-" + title,
					Title:       title,
					Link:        &feeds.Link{Href: "https://www.helios.pl" + url},
					Description: strings.Join(hours, "<br>") + "<br>(niepełna lista godzin i seansów)<br><br>Opis:<br>" + description,
					Created:     movieDate,
					Updated:     movieDate,
				})
			}
		})
	}
	return feed
}

func GetMoviesFromDay(day int) *goquery.Document {
	url := "https://www.helios.pl/23,Szczecin/Repertuar/axRepertoire/?dzien=" + strconv.Itoa(day) + "&kino=0"
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
	people1 := HeliosResponse{}
	jsonErr := json.Unmarshal(body, &people1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	node, err := html.Parse(strings.NewReader(people1.Html))
	if err != nil {
		log.Printf("error when parsing file: %v", err.Error())
	}

	return goquery.NewDocumentFromNode(node)
}
