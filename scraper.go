package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Cee-X/rssagg/internal/database"
	"github.com/google/uuid"
)



func startScrapping(db *database.Queries, concurency int, timeBetweenRequest time.Duration){
	log.Printf("Scrapping started with %d concurency and %s time between requests", concurency, timeBetweenRequest)
	ticket := time.NewTicker(timeBetweenRequest)
	for ; ; <- ticket.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurency),
		)
		if err != nil {
			log.Printf("Failed to get feeds to fetch: %v", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scarpeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scarpeFeed(db *database.Queries ,wg *sync.WaitGroup, feed database.Feed){
	defer wg.Done()
	log.Printf("Start scrapping feed %s", feed.Url)
	_, err := db.MarkFeedAsFetched(
		context.Background(),
		feed.ID,
	)
	if err != nil {
		log.Printf("Failed to mark feed as fetched: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Failed to fetch feed: %v", err)
		return
	}
	for _ , item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123, item.PubDate )
		if err != nil {
			log.Printf("Failed to parse time: %v", err)
			continue
		}
		_,err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url: item.Link,
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key"){
				continue
			}
			log.Printf("Failed to create post: %v", err)
		}
	}
	log.Printf("Feed %s collected, %d posts found", feed.Url, len(rssFeed.Channel.Item))
}