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

	channel := "example-channel"

	// Subscribe to the channel
	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()

	// Get the channel for receiving messages
	ch := sub.Channel()

	log.Printf("Subscribed to channel: %s", channel)

	// Listen for messages
	for msg := range ch {
		log.Printf("Received message: %s", msg.Payload)
	}
}
