package main

import (
	"context"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	ce "github.com/knative/pkg/cloudevents"
)

func newSinkPoster(sink string) (sinker *sinkPoster, err error) {

	// cloud events
	ceClient := ce.NewClient(sink, ce.Builder{
		EventType: "com.twitter",
		Source:    "https://api.twitter.com/1.1/search/tweets.json",
	})

	s := &sinkPoster{
		client: ceClient,
	}

	return s, nil

}

type sinkPoster struct {
	client *ce.Client
}

func (s *sinkPoster) post(ctx context.Context, t *twitter.Tweet) error {

	log.Printf("Posting tweet: %s\n", t.IDStr)
	eventTime, err := time.Parse(time.RubyDate, t.CreatedAt)

	if err != nil {
		log.Printf("Error while parsing created at: %v", err)
		eventTime = time.Now()
	}

	return s.client.Send(t, ce.V01EventContext{
		EventID:   t.IDStr,
		EventTime: eventTime,
	})

}
