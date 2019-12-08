package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/weriKK/dashboard/loggermw"
)

func webfeedListHandler(w http.ResponseWriter, r *http.Request) {
	//emit
	//wait for answer
}

func webfeedHandler(w http.ResponseWriter, r *http.Request) {
	//emit
	//wait for answer

	// the client sends separate request for each feed
	// these should be combined into one, and fan out with
	// messages in the api service. When the answers arrive
	// api service can combine them and return it to the client
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/webfeeds", loggermw.HandlerFunc(webfeedListHandler)).Methods("GET")
	r.HandleFunc("/webfeeds/{id:[0-9]+}", loggermw.HandlerFunc(webfeedHandler)).Methods("GET")

	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Infof("Listening on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
