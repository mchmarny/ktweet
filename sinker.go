package main

import (
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	ce "github.com/knative/pkg/cloudevents"
)

type sinker struct {
	client *ce.Client
}

func (s *sinker) post(t *twitter.Tweet) error {
	eventTime, err := time.Parse(time.RubyDate, t.CreatedAt)

	if err != nil {
		log.Printf("Error while parsing created at: %v", err)
		eventTime = time.Now()
	}

	return s.client.Send(t, ce.V02EventContext{
		ID:          t.IDStr,
		Time:        eventTime,
		ContentType: "application/json",
	})
}
