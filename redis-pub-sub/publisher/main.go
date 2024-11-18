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

	channel := "example-channel"

	for i := 1; i <= 10; i++ {
		message := fmt.Sprintf("Message %d", i)
		err := rdb.Publish(ctx, channel, message).Err()
		if err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}
		log.Printf("Published: %s", message)
		time.Sleep(1 * time.Second) // Simulate some delay between messages
	}

	log.Println("All messages published.")
}
