package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/config"
	"arena-matchmaker/pkg/logic"
	"arena-matchmaker/pkg/queue"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Arena Matchmaker...")

	cfg := config.Load()

	fmt.Printf("Connecting to Redis at: %s\n", cfg.RedisAddr)
	qClient, err := queue.NewClient(cfg.RedisAddr)

	if err != nil {
		log.Fatalf("Could not initialize Redis: %v", err)
	}
	defer qClient.RDB.Close()

	mmLogic := logic.NewMatchmaker(qClient.RDB, cfg.TickRate, cfg.MaxMMRDiff)

	ctx, cancel := context.WithCancel(context.Background())

	go mmLogic.Run(ctx)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	matchmaker.RegisterMatchmakerServiceServer(grpcServer, &GrpcServer{
		QueueClient: qClient,
	})

	go func() {
		fmt.Printf("🚀 gRPC Server listening on %s\n", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	fmt.Println("\n Shutting down server...")

	grpcServer.GracefulStop()
	fmt.Println("gRPC Server stopped.")

	cancel()

	time.Sleep(1 * time.Second)
	fmt.Println("Goodbye!")
}
