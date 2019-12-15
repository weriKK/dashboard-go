package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/weriKK/dashboard/config"
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

	feedSvcAPIRoot := config.GetString("FEED_SVC_API_ROOT")
	url := fmt.Sprintf("%s/feedidlist", feedSvcAPIRoot)

	resp, err := http.Get(url)
	if err != nil {
		return []FeedIdList{}, err
	}

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []FeedIdList{}, fmt.Errorf("http error received from feed service and can't read response body: %s - %s", resp.Status, err.Error())
		}

		errorMsg := fmt.Sprintf("http error received from feed service: %s - %s", resp.Status, string(msg))
		return []FeedIdList{}, fmt.Errorf(errorMsg)
	}

	var feedIDs []FeedIdList
	if err = json.NewDecoder(resp.Body).Decode(&feedIDs); err != nil {
		return []FeedIdList{}, err
	}

	return feedIDs, nil
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

	feedSvcAPIRoot := config.GetString("FEED_SVC_API_ROOT")
	url := fmt.Sprintf("%s/feedcontent/%d/%d", feedSvcAPIRoot, id, limit)

	resp, err := http.Get(url)
	if err != nil {
		return new(FeedContent), err
	}

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return new(FeedContent), fmt.Errorf("http error received from feed service and can't read response body: %s - %s", resp.Status, err.Error())
		}

		errorMsg := fmt.Sprintf("http error received from feed service: %s - %s", resp.Status, string(msg))
		return new(FeedContent), fmt.Errorf(errorMsg)
	}

	var feedContent FeedContent
	if err = json.NewDecoder(resp.Body).Decode(&feedContent); err != nil {
		return new(FeedContent), err
	}

	return &feedContent, nil
}

func webfeedListHandler(w http.ResponseWriter, r *http.Request) {
	feedList, err := getFeedIdList()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pubUrl := config.GetString("SVC_PUBLIC_API_ROOT")

	payload := jsonFeedList{len(feedList), []jsonFeedListItem{}}
	for _, v := range feedList {
		payload.Feeds = append(payload.Feeds, jsonFeedListItem{v.Name, v.Url, fmt.Sprintf("%s%s/%d", pubUrl, r.URL, v.Id), v.Column, v.ItemLimit})
	}
	writeJSONPayload(w, payload)
}

func webFeedContentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// if error, limit is set to 0
	limit, err := GetLimitQueryParam(r.URL.Query())
	if err != nil {
		limit = 10
	}

	feedContent, err := getFeedContent(id, limit)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler).Methods("GET")
	r.HandleFunc("/webfeeds", loggermw.HandlerFunc(webfeedListHandler)).Methods("GET")
	r.HandleFunc("/webfeeds/{id:[0-9]+}", loggermw.HandlerFunc(webFeedContentHandler)).Methods("GET")

	s := &http.Server{
		Addr:    ":" + config.GetString("SVC_PORT"),
		Handler: r,
	}

	log.Infof("Listening on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
