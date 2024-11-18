package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	streamKey := "example-stream"
	lastID := "0" // Start reading from the beginning of the stream

	log.Printf("Consumer is listening to stream: %s", streamKey)

	for {
		// Read messages from the stream
		messages, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{streamKey, lastID}, // Read from the last ID
			Block:   0,                           // Wait indefinitely for new messages
		}).Result()
		if err != nil {
			log.Fatalf("Failed to read from stream: %v", err)
		}

		for _, stream := range messages {
			for _, msg := range stream.Messages {
				log.Printf("Received message: %v", msg.Values["message"])
				lastID = msg.ID // Update lastID to continue from the latest message
			}
		}
	}
}
