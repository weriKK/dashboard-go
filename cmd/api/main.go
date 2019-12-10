package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/weriKK/dashboard/loggermw"
)

type jsonFeedList struct {
	Count int
	Feeds []jsonFeedListItem
}

type jsonFeedListItem struct {
	Name      string
	Url       string
	Resource  string
	Column    int
	ItemLimit int
}

type jsonFeed struct {
	Id    int
	Name  string
	Url   string
	Items []jsonFeedItem
}

type jsonFeedItem struct {
	Title       string
	Url         string
	Description string
	Published   string
}

type FeedIdList struct {
	Id        int
	Name      string
	Url       string
	Rss       string
	Column    int
	ItemLimit int
}

func getFeedIdList() ([]FeedIdList, error) {
	return []FeedIdList{}, nil
}

type FeedItem struct {
	Title       string
	Url         string
	Description string
	Published   string
}

type FeedContent struct {
	Id     int
	Name   string
	Url    string
	Column int
	Items  []FeedItem
}

func getFeedContent(id int, limit int) (*FeedContent, error) {
	return new(FeedContent), nil
}

func webfeedListHandler(w http.ResponseWriter, r *http.Request) {
	feedList, _ := getFeedIdList()

	payload := jsonFeedList{len(feedList), []jsonFeedListItem{}}
	for _, v := range feedList {
		payload.Feeds = append(payload.Feeds, jsonFeedListItem{v.Name, v.Url, fmt.Sprintf("http://%s%s/%d", r.Host, r.URL, v.Id), v.Column, v.ItemLimit})
	}
	writeJSONPayload(w, payload)
}

func webFeedContentHandler(w http.ResponseWriter, r *http.Request) {
	parsed, err := ParseUrl(r.URL)
	if err != nil {
		log.Panic(err)
	}

	id, err := strconv.Atoi(parsed.LastPath)
	if err != nil {
		log.Panic(err)
	}

	// if error, limit is set to 0
	limit, err := parsed.GetLimitQueryParam()

	feedContent, _ := getFeedContent(id, limit)

	payload := jsonFeed{feedContent.Id, feedContent.Name, feedContent.Url, []jsonFeedItem{}}
	for _, v := range feedContent.Items {
		payload.Items = append(payload.Items, jsonFeedItem{v.Title, v.Url, v.Description, v.Published})
	}
	writeJSONPayload(w, payload)
}

func writeJSONPayload(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/webfeeds", loggermw.HandlerFunc(webfeedListHandler)).Methods("GET")
	r.HandleFunc("/webfeeds/{id:[0-9]+}", loggermw.HandlerFunc(webFeedContentHandler)).Methods("GET")

	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Infof("Listening on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
