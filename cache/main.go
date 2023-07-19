package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Create a Redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a store.
	// store := NewStore(client)

	// Channel to receive broadcast messages from the "cache_channel".
	broadcastChannel := make(chan []*FileChange)
	var wg sync.WaitGroup
	wg.Add(1)

	// Start subscribing to the Redis channel "cache_channel".
	go subscribeToChannel(ctx, client, "cache_channel", broadcastChannel)

	// Start a goroutine to handle WebSocket connections.
	go handleWebSocketConnections()

	// Wait for termination signal to gracefully shut down the service.
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	<-terminate
	cancel()  // Cancel the context to stop the goroutines.
	wg.Wait() // Wait for the broadcast goroutine to finish gracefully.
}
