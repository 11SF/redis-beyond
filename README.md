
# Redis Beyond Basics

## HyperLogLog
Run the example to demonstrate unique count estimation using Redis HyperLogLog:

```bash
go run main.go
```

---

## Redis Pub/Sub

### Publisher
Start the publisher to send messages to a channel:

```bash
go run main.go
```

### Subscriber
Start a subscriber to receive messages from the channel:

```bash
go run main.go
```

Simulate processing delay for the subscriber:

```bash
go run main.go -delay 2000
```

---

## Redis Streams

### Producer
Run the producer to add messages to a stream:

```bash
go run main.go
```

### Consumers
1. Join a consumer group and process messages:
   ```bash
   go run main.go -name consumer1 -group true -group-name example-group -stream example-stream
   ```
   ```bash
   go run main.go -name consumer2 -group true -group-name example-group -stream example-stream
   ```

2. Process all messages independently (without a group):
   ```bash
   go run main.go -name consumer3 -stream example-stream
   ```

---

## Redis Pipeline
Run the example to benchmark performance with and without a pipeline:

```bash
go run main.go
```

---

## Redis Transactions
Run the example to simulate atomic operations with Redis transactions:

```bash
go run main.go
```
