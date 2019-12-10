package main

import (
	"errors"
	"sort"
	"sync"
)

type FeedRepository interface {
	GetAll(limit int) ([]Feed, error)
	GetById(id int) (*Feed, error)
	Count() (int, error)
}

type entry struct {
	name   string
	url    string
	rss    string
	column int
	limit  int
}

type inMemoryFeedRepository struct {
	data map[int]*entry
	mux  sync.RWMutex
}

func (db *inMemoryFeedRepository) GetAll(limit int) ([]Feed, error) {

	var data []Feed
	count := 0

	for k, v := range db.data {
		if limit <= count {
			break
		}
		newFeed := NewFeed(k, v.name, v.url, v.rss, v.column, v.limit)
		data = append(data, newFeed)
		count++
	}

	// Fun fact, map stores elements in a randomized order,
	// can't rely on it keeping indices in the expected order
	sort.Slice(data, func(i, j int) bool {
		return data[i].Id() < data[j].Id()
	})

	return data, nil
}

func (db *inMemoryFeedRepository) GetById(id int) (*Feed, error) {

	newFeed := Feed{}

	for k := range db.data {
		if k == id {
			newFeed = NewFeed(id, db.data[id].name, db.data[id].url, db.data[id].rss, db.data[id].column, db.data[id].limit)
			return &newFeed, nil
		}
	}

	return &newFeed, errors.New("Feed with given id not found!")

}

func (db *inMemoryFeedRepository) Count() (int, error) {
	return len(db.data), nil
}

func (db *inMemoryFeedRepository) add(value *entry) {
	db.mux.Lock()
	db.data[len(db.data)] = value
	db.mux.Unlock()
}

func (db *inMemoryFeedRepository) initializeWithData() {
	db.add(&entry{"MMO-Champion", "https://www.mmo-champion.com", "http://www.mmo-champion.com/external.php?do=rss&type=newcontent&sectionid=1&days=120&count=20", 0, 5})
	db.add(&entry{"Reddit - Games", "https://www.reddit.com/r/Games/", "https://www.reddit.com/r/Games/.rss", 0, 12})
	db.add(&entry{"Programming Praxis", "https://programmingpraxis.com", "https://programmingpraxis.com/feed/", 0, 5})
	db.add(&entry{"Handmade Hero", "https://programmingpraxis.com", "https://www.youtube.com/feeds/videos.xml?channel_id=UCaTznQhurW5AaiYPbhEA-KA", 0, 4})
	db.add(&entry{"GiantBomb", "https://www.giantbomb.com", "http://www.giantbomb.com/feeds/mashup/", 1, 6})
	db.add(&entry{"RockPaperShotgun", "https://www.rockpapershotgun.com", "http://feeds.feedburner.com/RockPaperShotgun", 1, 6})
	db.add(&entry{"Shacknews", "https://www.shacknews.com", "http://www.shacknews.com/rss?recent_articles=1", 1, 8})
	db.add(&entry{"Bluenews", "https://www.bluesnews.com", "http://www.bluesnews.com/news/news_1_0.rdf", 1, 8})
	db.add(&entry{"Gamasutra", "https://www.gamasutra.com/", "http://feeds.feedburner.com/GamasutraFeatureArticles/", 2, 4})
	db.add(&entry{"ArsTechnica", "https://arstechnica.com", "http://feeds.arstechnica.com/arstechnica/index", 2, 6})
	db.add(&entry{"GamesIndustry", "https://www.gamesindustry.biz", "http://www.gamesindustry.biz/rss/gamesindustry_news_feed.rss", 2, 6})
	db.add(&entry{"Y Combinator", "https://news.ycombinator.com", "https://news.ycombinator.com/rss", 2, 12})
}

func NewInMemoryFeedRepository() *inMemoryFeedRepository {
	db := inMemoryFeedRepository{data: make(map[int]*entry)}
	db.initializeWithData()
	return &db
}
