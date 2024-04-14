package queue

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
)

type MemoryQueueAdapter struct{}

func NewMemoryQueueAdapter() *MemoryQueueAdapter {
	return &MemoryQueueAdapter{}
}

func (eb *MemoryQueueAdapter) Publish(ctx context.Context, eventPayload interface{}) error {
	eventType := reflect.TypeOf(eventPayload)
	payloadJson, _ := json.Marshal(eventPayload)
	log.Printf("** [Publish] %s: %v ---", eventType, string(payloadJson))
	return nil
}
