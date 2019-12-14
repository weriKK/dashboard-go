package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/mmcdole/gofeed"
	"github.com/weriKK/dashboard/config"
	"github.com/weriKK/dashboard/loggermw"
)

type Format int

type Feed struct {
	id               int
	name             string
	url              string
	rss              string
	column           int
	visibleItemLimit int
}

var appdb FeedRepository

func init() {
	appdb = NewInMemoryFeedRepository()
}

func NewFeed(id int, name string, url string, rss string, column int, visibleItemLimit int) Feed {
	return Feed{id, name, url, rss, column, visibleItemLimit}
}

func (f Feed) Id() int      { return f.id }
func (f Feed) Name() string { return f.name }
func (f Feed) Url() string  { return f.url }
func (f Feed) Rss() string  { return f.rss }
func (f Feed) Column() int  { return f.column }
func (f Feed) Limit() int   { return f.visibleItemLimit }

type FeedIdList struct {
	Id        int
	Name      string
	Url       string
	Rss       string
	Column    int
	ItemLimit int
}

func webfeedIdListHandler(w http.ResponseWriter, r *http.Request) {

	count, _ := appdb.Count()
	feeds, _ := appdb.GetAll(count)

	feedInfo := []FeedIdList{}

	for _, v := range feeds {
		feedInfo = append(feedInfo, FeedIdList{v.Id(), v.Name(), v.Url(), v.Rss(), v.Column(), v.Limit()})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feedInfo); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func webFeedContentHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	limit, _ := strconv.Atoi(mux.Vars(r)["limit"])

	feed, _ := appdb.GetById(id)

	feedParser := gofeed.NewParser()
	parsed, err := feedParser.ParseURL(feed.Rss())
	if err != nil {
		log.Errorf("Failed to parse RSS feed '%v': %v", feed.Rss(), err)
		http.Error(w, "Failed to parse RSS feed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if limit == 0 {
		limit = len(parsed.Items)
	}

	limit = min(limit, len(parsed.Items))

	items := []FeedItem{}
	for itemIdx := 0; itemIdx < limit; itemIdx++ {
		p := parsed.Items[itemIdx]
		items = append(items, FeedItem{p.Title, p.Link, p.Description, p.Published})
	}

	content := FeedContent{feed.Id(), parsed.Title, feed.Url(), feed.Column(), items}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(content); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler).Methods("GET")
	r.HandleFunc("/feedidlist", loggermw.HandlerFunc(webfeedIdListHandler)).Methods("GET")
	r.HandleFunc("/feedcontent/{id:[0-9]+}/{limit:[0-9]+}", loggermw.HandlerFunc(webFeedContentHandler)).Methods("GET")

	s := &http.Server{
		Addr:    ":" + config.GetString("SVC_PORT"),
		Handler: r,
	}

	log.Infof("Listening on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
