package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store struct {
	client *redis.Client
	mutex  *sync.Mutex
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
		mutex:  &sync.Mutex{},
	}
}

func (s *Store) Set(ctx context.Context, key string, value string, exp int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	expiration := time.Duration(exp) * time.Second
	err := s.client.Set(ctx, key, value, expiration)
	if err != nil {
		log.Println("Error setting value in Redis:", err)
	}
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("Error getting value from Redis:", err)
		}
		return "", err
	}

	return val, nil
}

// Message structure to represent file changes sent from the Cache Service.
type FileChange struct {
	SessionID string // Assuming each session has a unique ID.
	UserID    string // Assuming each user has a unique ID.
	FileID    string // Assuming each file has a unique ID.
	Content   string // New content of the file.
}

// Function to process file change messages and prepare data to broadcast to clients via WebSocket.
func processFileChanges(message string) []*FileChange {
	// Here, you can implement the logic to parse and process the message payload.
	// For this example, I'm assuming the payload is a JSON string representing an array of FileChange objects.
	var fileChanges []*FileChange
	err := json.Unmarshal([]byte(message), &fileChanges)
	if err != nil {
		log.Println("Error parsing message payload:", err)
		return nil
	}
	return fileChanges
}

// Subscribe to Redis channel and process incoming messages.
func subscribeToChannel(ctx context.Context, client *redis.Client, channel string, broadcastChannel chan<- []*FileChange) {
	defer close(broadcastChannel)

	sub := client.Subscribe(ctx, channel)
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		fileChanges := processFileChanges(msg.Payload)
		if fileChanges != nil {
			broadcastChannel <- fileChanges
		}
	}
}

// Function to publish a message to "cache_channel".
func publishMessage(ctx context.Context, client *redis.Client, channel, message string) error {
	err := client.Publish(ctx, channel, message).Err()
	return err
}
