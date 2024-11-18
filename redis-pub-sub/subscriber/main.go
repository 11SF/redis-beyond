package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Define command-line flag for delay
	delay := flag.Int("delay", 0, "Delay in milliseconds to simulate message processing time")
	flag.Parse()

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

		// Simulate processing delay
		if *delay > 0 {
			time.Sleep(time.Duration(*delay) * time.Millisecond)
		}
	}
}
