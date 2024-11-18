package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	streamKey := "example-stream-group"
	groupName := "example-group"
	consumerName := "consumer-1" // Change this for multiple consumers

	// Create the consumer group if it doesn't already exist
	err := rdb.XGroupCreateMkStream(ctx, streamKey, groupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	log.Printf("Consumer %s listening on group %s, stream %s", consumerName, groupName, streamKey)

	for {
		// Read messages from the stream for this consumer
		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamKey, ">"},
			Block:    0, // Wait indefinitely for new messages
		}).Result()
		if err != nil {
			log.Fatalf("Failed to read from stream: %v", err)
		}

		for _, stream := range messages {
			for _, msg := range stream.Messages {
				log.Printf("Consumer %s received message: %v", consumerName, msg.Values["message"])

				// Acknowledge the message
				err := rdb.XAck(ctx, streamKey, groupName, msg.ID).Err()
				if err != nil {
					log.Printf("Failed to acknowledge message: %v", err)
				}
			}
		}
	}
}
