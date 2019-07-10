package main

import (
	"context"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var (
	consumerKey    = mustGetEnvVar("T_CONSUMER_KEY", "")
	consumerSecret = mustGetEnvVar("T_CONSUMER_SECRET", "")
	accessToken    = mustGetEnvVar("T_ACCESS_TOKEN", "")
	accessSecret   = mustGetEnvVar("T_ACCESS_SECRET", "")
)

func search(ctx context.Context, query, sink string, stop <-chan struct{}) {

	// twitter client config
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	twClient := twitter.NewClient(httpClient)

	sinker, err := newSinkPoster(sink)
	if err != nil {
		log.Fatalf("Error getting sinker: %v", err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(t *twitter.Tweet) {
		log.Printf("Got tweet: %s\n", t.IDStr)
		if err := sinker.post(ctx, t); err != nil {
			log.Printf("Error on tweet handle: %v\n", err)
		}
	}

	params := &twitter.StreamFilterParams{
		Track:         []string{query},
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	stream, err := twClient.Streams.Filter(params)
	if err != nil {
		log.Fatalf("Error on filter create: %v\n", err)
		return
	}

	log.Printf("Starting tweet streamming for: %s\n", query)
	go demux.HandleChan(stream.Messages)

}

func mustGetEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	if fallbackValue == "" {
		panic("Required env var not set: " + key)
	}

	return fallbackValue
}
