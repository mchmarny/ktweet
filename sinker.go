package main

import (
	"context"
	"log"
	"time"

	ce "github.com/cloudevents/sdk-go"
	"github.com/dghubble/go-twitter/twitter"
)

func newSinkPoster(targetURI string) (sinker *sinkPoster, err error) {

	// cloud events
	ceTransport, err := ce.NewHTTPTransport(
		ce.WithTarget(targetURI),
		ce.WithEncoding(ce.HTTPBinaryV02),
	)
	if err != nil {
		log.Fatalf("Error while creating CloudEvents transport: %v\n", err)
	}
	ceClient, err := ce.NewClient(ceTransport)
	if err != nil {
		log.Fatalf("Error while creating CloudEvents client: %v\n", err)
	}

	s := &sinkPoster{
		client: ceClient,
	}

	return s, nil

}

type sinkPoster struct {
	client ce.Client
}

func (s *sinkPoster) post(ctx context.Context, t *twitter.Tweet) error {

	log.Printf("Posting tweet: %s\n", t.IDStr)
	eventTime, err := time.Parse(time.RubyDate, t.CreatedAt)

	if err != nil {
		log.Printf("Error while parsing created at: %v", err)
		eventTime = time.Now()
	}

	e := ce.NewEvent("0.2")
	e.SetSpecVersion("0.2")
	e.SetID(t.IDStr)
	e.SetTime(eventTime)
	e.SetType("com.twitter")
	e.SetSource("https://api.twitter.com/1.1/search/tweets.json")
	e.SetDataContentType("application/json")
	e.SetData(t)

	result, err := s.client.Send(ctx, e)

	if result != nil {
		log.Printf("Result: %v", result)
	}

	return err

}
