package main

import (
	"context"
	"flag"
	"log"

	"github.com/knative/pkg/signals"
)

var (
	sink  string
	query string
)

func init() {
	flag.StringVar(&sink, "sink", "", "where to sink events to")
	flag.StringVar(&query, "query", "", "twitter query/search string")
}

func main() {

	flag.Parse()

	ctx := context.Background()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	if query == "" {
		log.Fatal("Query parameter required")
	}

	log.Printf("Start (sink: %s, query: %s)", sink, query)

	search(ctx, query, sink, stopCh)

	<-stopCh
}
