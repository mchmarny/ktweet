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
}

type searcher struct {
	client  *twitter.Client
	query   string
	handler func(*twitter.Tweet) error
	stop    <-chan struct{}
}

func (s *searcher) run() {

	params := &twitter.StreamFilterParams{
		Track:         []string{s.query},
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(t *twitter.Tweet) {
		log.Printf("Got tweet from %s\n", t.User.Name)
		if handleErr := s.handler(t); handleErr != nil {
			log.Printf("Failed to post: %v\n", handleErr)
		}
	}

	stream, err := s.client.Streams.Filter(params)
	if err != nil {
		log.Printf("Failed to create filter: %v\n", err)
		return
	}
	go demux.HandleChan(stream.Messages)
}
