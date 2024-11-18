package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	streamKey := "example-stream-group"

	for i := 1; i <= 10; i++ {
		message := fmt.Sprintf("Message %d", i)
		_, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamKey,
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
