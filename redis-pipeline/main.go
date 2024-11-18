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
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Number of keys to insert
	totalKeys := 100000

	// Without Pipeline
	start := time.Now()
	for i := 0; i < totalKeys; i++ {
		key := fmt.Sprintf("user:%d:name", i)
		value := fmt.Sprintf("User %d", i)
		err := rdb.Set(ctx, key, value, 5*time.Minute).Err()
		if err != nil {
			log.Fatalf("Error setting key without pipeline: %v", err)
		}
	}
	durationWithoutPipeline := time.Since(start)
	fmt.Printf("Time taken without pipeline: %v\n", durationWithoutPipeline)

	// With Pipeline
	start = time.Now()
	pipe := rdb.Pipeline()
	for i := 0; i < totalKeys; i++ {
		key := fmt.Sprintf("user_pipeline:%d:name", i)
		value := fmt.Sprintf("User %d", i)
		pipe.Set(ctx, key, value, 5*time.Minute)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing pipeline: %v", err)
	}
	durationWithPipeline := time.Since(start)
	fmt.Printf("Time taken with pipeline: %v\n", durationWithPipeline)
}
