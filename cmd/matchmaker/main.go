package main

import (
	"arena-matchmaker/pkg/config"
	"arena-matchmaker/pkg/logic"
	"arena-matchmaker/pkg/queue"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting Arena Matchmaker...")

	cfg := config.Load()

	fmt.Printf("🔌 Connecting to Redis at: %s\n", cfg.RedisAddr)
	qClient, err := queue.NewClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Could not initialize Redis: %v", err)
	}
	defer qClient.RDB.Close()

	mmLogic := logic.NewMatchmaker(qClient.RDB, cfg.TickRate, cfg.MaxMMRDiff)

	ctx, cancel := context.WithCancel(context.Background())
	go mmLogic.Run(ctx)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	fmt.Println("\n⚠️  Shutting down Matchmaker...")
	cancel()
	fmt.Println("✅ Matchmaker stopped.")
}
