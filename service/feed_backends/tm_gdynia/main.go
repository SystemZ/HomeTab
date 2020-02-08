package tm_gdynia

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
		Title:       "TM Gdynia",
		Link:        &feeds.Link{Href: "https://www.muzyczny.org/pl/repertuar.html"},
		Description: "Latest plays",
	}
	node := feed_backends.GetDocOnline("https://www.muzyczny.org/pl/repertuar.html")
	//node := common.GetDocFromDisk("index.html")
	doc := goquery.NewDocumentFromNode(node)
	doc.Find(".spektakl_row").Each(func(i int, s *goquery.Selection) {
		playYear, _ := s.Attr("data-year")
		playMonth, _ := s.Attr("data-month")
		playDay, _ := s.Attr("data-day")
		playTitle := s.Find(".spektakl_szczegoly").Find(".h2").Text()
		playDatePlace := s.Find(".spektakl_szczegoly").Find(".h6").Text()

		// parse date of play
		playDatePlaceSplit := strings.Split(strings.TrimSpace(playDatePlace), "/")
		playDateTime := playDatePlaceSplit[0]
		playDateTimeSplit := strings.Split(playDateTime, ".")
		playDateHour := playDateTimeSplit[0]
		playDateMinute := strings.TrimSpace(playDateTimeSplit[1])
		playTime := playDateHour + ":" + playDateMinute
		layout := "2006-01-02 15:04:05 -0700"
		dateString := playYear + "-" + playMonth + "-" + playDay + " " + playTime + ":00 +0200"
		t, err := time.Parse(layout, dateString)
		if err != nil {
			log.Println("Error while parsing date :", err)
		}

		// just a little formatting...
		feedTitle := playTitle + " " + playTime + " " + playYear + "-" + playMonth + "-" + playDay

		// ... and bang!
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          feedTitle,
			Title:       feedTitle,
			Link:        &feeds.Link{Href: "https://www.muzyczny.org/pl/repertuar.html"},
			Description: playDatePlace,
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
