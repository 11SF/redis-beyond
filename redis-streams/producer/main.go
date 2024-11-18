package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Define command-line flags
	streamKey := flag.String("stream", "example-stream", "Redis stream key")
	flag.Parse()

	// Validate required flags
	if *streamKey == "" {
		log.Fatal("Stream key cannot be empty")
	}

	// Initialize Redis client
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	log.Printf("Producing messages to stream: %s", *streamKey)

	// Publish messages to the stream
	for i := 1; i <= 10; i++ {
		message := fmt.Sprintf("Message %d", i)
		_, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: *streamKey,
			Values: map[string]interface{}{"message": message},
		}).Result()
		if err != nil {
			log.Fatalf("Failed to add message to stream: %v", err)
		}
		log.Printf("Published: %s", message)
		time.Sleep(1 * time.Second) // Simulate delay between messages
	}

	log.Println("All messages published.")
}
