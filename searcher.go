package main

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

func runSearcher(client *twitter.Client,
	query string,
	handler func(*twitter.Tweet) error,
	stop <-chan struct{}) {
	searcher := &searcher{
		client:  client,
		query:   query,
		handler: handler,
		stop:    stop,
	}
	searcher.run()
	log.Println("Twitter stream started...")
}

type searcher struct {
	client  *twitter.Client
	query   string
	handler func(*twitter.Tweet) error
	stop    <-chan struct{}
}

func (s *searcher) run() {

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(t *twitter.Tweet) {
		log.Printf("Got tweet: %s\n", t.IDStr)
		if err := s.handler(t); err != nil {
			log.Printf("Error on tweet handle: %v\n", err)
		}
	}

	params := &twitter.StreamFilterParams{
		Track:         []string{s.query},
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	stream, err := s.client.Streams.Filter(params)
	if err != nil {
		log.Fatalf("Error on filter create: %v\n", err)
		return
	}

	log.Printf("Starting tweet streamming for: %s\n", s.query)
	go demux.HandleChan(stream.Messages)
}
