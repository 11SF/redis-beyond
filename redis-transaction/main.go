package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func transferFunds(ctx context.Context, wg *sync.WaitGroup, rdb *redis.Client, fromAccount, toAccount string, amount int64) error {

	defer wg.Done()
	slog.Info("Starting transfer", slog.String("from", fromAccount), slog.String("to", toAccount), slog.Int64("amount", amount))

	// Watch the keys to detect changes
	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		fromBalance, err := tx.Get(ctx, fromAccount).Int64()
		if err != nil {
			return fmt.Errorf("could not get balance for %s: %v", fromAccount, err)
		}

		if fromBalance < amount {
			return fmt.Errorf("insufficient funds in %s", fromAccount)
		}

		// Simulate processing delay
		time.Sleep(5 * time.Second)

		_, err = tx.Get(ctx, toAccount).Int64()
		if err != nil {
			return fmt.Errorf("could not get balance for %s: %v", fromAccount, err)
		}

		// Begin the transaction
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.DecrBy(ctx, fromAccount, amount)
			pipe.IncrBy(ctx, toAccount, amount)
			return nil
		})
		return err
	}, fromAccount, toAccount)

	if err != nil {
		slog.Error("Transaction failed", slog.String("err", err.Error()))
		return err
	}
	slog.Info("Transaction executed successfully")
	return nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Set initial balances
	slog.Info("Initializing account A with balance 500")
	rdb.Set(ctx, "account:A", 500, 5*time.Minute)
	slog.Info("Initializing account B with balance 500")
	rdb.Set(ctx, "account:B", 500, 5*time.Minute)

	// Perform fund transfer in a separate goroutine
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go transferFunds(ctx, wg, rdb, "account:A", "account:B", 200)

	// Simulate an update to account:B during the transfer
	if true {
		time.Sleep(2 * time.Second)
		slog.Info("Simulating external update: Setting account:B balance to 1000")
		err := rdb.Set(ctx, "account:B", 1000, 5*time.Minute).Err()
		if err != nil {
			slog.Error("Failed to update account:B", slog.String("err", err.Error()))
		}
	}

	wg.Wait()
	slog.Info("All operations completed.")

	balanceA, _ := rdb.Get(ctx, "account:A").Result()
	balanceB, _ := rdb.Get(ctx, "account:B").Result()

	slog.Info("Final balance of account A", slog.String("balance", balanceA))
	slog.Info("Final balance of account B", slog.String("balance", balanceB))
}
