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

	log.Printf("Posting tweet: %s\n", t.IDStr)
	eventTime, err := time.Parse(time.RubyDate, t.CreatedAt)

	if err != nil {
		log.Printf("Error while parsing created at: %v", err)
		eventTime = time.Now()
	}

	err = s.client.Send(t, ce.V02EventContext{
		ID:          t.IDStr,
		Time:        eventTime,
		ContentType: "application/json",
	})

	if err != nil {
		log.Printf("Error posting: %v", err)
	}

	return err

}
