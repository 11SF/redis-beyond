package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Define the HyperLogLog key
	hyperLogLogKey := "test-hyperloglog"

	// Clean up before starting
	rdb.Del(ctx, hyperLogLogKey)

	// Number of elements to add
	nExpected := 1000

	log.Printf("Adding %d unique elements to HyperLogLog...", nExpected)

	// Add elements to the HyperLogLog
	for i := 1; i <= nExpected; i++ {
		element := fmt.Sprintf("user%d", i)
		err := rdb.PFAdd(ctx, hyperLogLogKey, element).Err()
		if err != nil {
			log.Fatalf("Failed to add element to HyperLogLog: %v", err)
		}
	}

	// Get the approximate cardinality (unique count)
	Approximate, err := rdb.PFCount(ctx, hyperLogLogKey).Result()
	if err != nil {
		log.Fatalf("Failed to get HyperLogLog count: %v", err)
	}

	// Print results
	fmt.Printf("\nExpected Count (nExpected): %d\n", nExpected)
	fmt.Printf("Approximate Count: %d\n", Approximate)
	fmt.Printf("Difference: %d\n", Approximate-int64(nExpected))

	// Cleanup
	err = rdb.Del(ctx, hyperLogLogKey).Err()
	if err != nil {
		log.Fatalf("Failed to delete HyperLogLog key: %v", err)
	}
	log.Println("\nHyperLogLog test completed and key deleted.")
}
