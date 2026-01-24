package server

import (
	"arena-matchmaker/pkg/queue"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting Arena Matchmaker...")

	redisClient, err := queue.NewClient("localhost:6379")
	if err != nil {
		log.Fatalf("Could not initialize Redis: %v", err)
	}
	defer redisClient.Close()

	fmt.Println("✅ Successfully connected to Redis!")

	// Keep the app running
	// for now, we just exit successfully.
}
