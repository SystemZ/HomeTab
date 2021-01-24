package generator

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
)

func GetDocOnline(url string) *html.Node {
	// https://stackoverflow.com/questions/13263492/set-useragent-in-http-request/13263993#13263993
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Println(string(body))

	node, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("error when parsing file: %v", err.Error())
	}
	return node
}

func GetDocFromDisk(filename string) *html.Node {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error when parsing file %v :%v", filename, err.Error())
	}
	file := bytes.NewReader(b)
	node, err := html.Parse(file)
	if err != nil {
		log.Printf("error when parsing file: %v", err.Error())
	}
	//log.Printf("%v", node.FirstChild.NextSibling.Data)
	return node
}

func ExampleScrape(node *html.Node) {
	// Load the HTML document
	doc := goquery.NewDocumentFromNode(node)
	// Find the review items
	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}
