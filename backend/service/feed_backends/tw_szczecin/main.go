package tw_szczecin

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"gitlab.com/systemz/tasktab/service/feed_backends"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func Fetch() *feeds.Feed {
	feed := &feeds.Feed{
		Title:       "TW Szczecin",
		Link:        &feeds.Link{Href: "http://wspolczesny.szczecin.pl"},
		Description: "Latest plays",
	}

	node := feed_backends.GetDocOnline("http://wspolczesny.szczecin.pl/repertuar/")
	//node := common.GetDocFromDisk("index.html")
	doc := goquery.NewDocumentFromNode(node)
	doc.Find(".repertuarTw-row").Each(func(i int, s *goquery.Selection) {
		playTitle := s.Find("a").First().Text()
		playDate := s.Find(".repertuarTw-dates__day-month").Text()
		playTime := s.Find(".repertuarTw-dates__time").Text()
		playPlace := s.Find(".repertuarTw-play-info__scene").Text()

		// parse date of play
		playDateSplit := strings.Split(playDate, ".")
		layout := "2006-01-02 15:04:05 -0700"
		dateString := "2019-" + playDateSplit[1] + "-" + playDateSplit[0] + " " + playTime + ":00 +0200"
		t, err := time.Parse(layout, dateString)
		if err != nil {
			log.Println("Error while parsing date :", err)
		}

		// just a little formatting...
		feedTitle := playTitle + " " + playTime + " " + playDate + " 2019"

		// ... and bang!
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          feedTitle,
			Title:       feedTitle,
			Link:        &feeds.Link{Href: "http://wspolczesny.szczecin.pl/repertuar/"},
			Description: playPlace,
			Created:     t,
			Updated:     t,
		})
	})

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
