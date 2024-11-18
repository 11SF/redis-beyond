package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Define command-line flags
	consumerName := flag.String("name", "default-consumer", "Name of the consumer")
	joinGroup := flag.Bool("group", false, "Join consumer group (true/false)")
	groupName := flag.String("group-name", "example-group", "Name of the consumer group")
	streamKey := flag.String("stream", "example-stream", "Redis stream key")
	flag.Parse()

	// Validate required flags
	if *consumerName == "" {
		fmt.Println("Consumer name cannot be empty")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	log.Printf("Consumer Name: %s, Join Group: %v", *consumerName, *joinGroup)

	if *joinGroup {
		// Join a consumer group
		err := rdb.XGroupCreateMkStream(ctx, *streamKey, *groupName, "$").Err()
		if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
			log.Fatalf("Failed to create consumer group: %v", err)
		}
		log.Printf("Joined consumer group: %s", *groupName)

		// Read messages using XREADGROUP
		for {
			messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    *groupName,
				Consumer: *consumerName,
				Streams:  []string{*streamKey, ">"},
				Block:    0, // Wait indefinitely for new messages
			}).Result()
			if err != nil {
				log.Fatalf("Failed to read from stream using consumer group: %v", err)
			}

			for _, stream := range messages {
				for _, msg := range stream.Messages {
					log.Printf("[%s] Received message: %v", *consumerName, msg.Values)
					// Acknowledge the message
					err := rdb.XAck(ctx, *streamKey, *groupName, msg.ID).Err()
					if err != nil {
						log.Printf("[%s] Failed to acknowledge message: %v", *consumerName, err)
					}
				}
			}
		}
	} else {
		// Read messages independently using XREAD
		lastID := "0" // Start from the beginning of the stream
		for {
			messages, err := rdb.XRead(ctx, &redis.XReadArgs{
				Streams: []string{*streamKey, lastID},
				Block:   0, // Wait indefinitely for new messages
			}).Result()
			if err != nil {
				log.Fatalf("Failed to read from stream: %v", err)
			}

			for _, stream := range messages {
				for _, msg := range stream.Messages {
					log.Printf("[%s] Received message: %v", *consumerName, msg.Values)
					lastID = msg.ID // Update the lastID to continue from the latest message
				}
			}
		}
	}
}
