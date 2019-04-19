package main

import (
	"flag"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	ce "github.com/knative/pkg/cloudevents"
	"github.com/knative/pkg/signals"
)

var (
	consumerKey    = mustGetEnvVar("T_CONSUMER_KEY", "")
	consumerSecret = mustGetEnvVar("T_CONSUMER_SECRET", "")
	accessToken    = mustGetEnvVar("T_ACCESS_TOKEN", "")
	accessSecret   = mustGetEnvVar("T_ACCESS_SECRET", "")

	sink  string
	query string
)

func init() {
	flag.StringVar(&sink, "sink", "", "where to sink events to")
	flag.StringVar(&query, "query", "", "twitter query/search string")
}

func main() {

	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	if query == "" {
		log.Fatal("Query string required")
	}

	log.Printf("Starting with sink: %s, query: %s", sink, query)

	ceClient := ce.NewClient(sink, ce.Builder{
		EventType:        "com.twitter",
		Source:           "https://api.twitter.com/1.1/search/tweets.json",
		EventTypeVersion: "0.2",
		SchemaURL:        "https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/tweet-object.html",
	})

	s := &sinker{client: ceClient}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	twClient := twitter.NewClient(httpClient)

	runSearcher(twClient, query, s.post, stopCh)

	<-stopCh
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
