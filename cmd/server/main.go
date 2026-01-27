package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/logic"
	"arena-matchmaker/pkg/queue"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Arena Matchmaker...")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	fmt.Printf("Connecting to Redis at: %s\n", redisAddr)

	qClient, err := queue.NewClient(redisAddr)
	if err != nil {
		log.Fatalf("Could not initialize Redis: %v", err)
	}
	defer qClient.RDB.Close()
	fmt.Println("Connected to Redis")

	mmLogic := logic.NewMatchmaker(qClient.RDB)

	go mmLogic.Run(context.Background())

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	matchmaker.RegisterMatchmakerServiceServer(grpcServer, &GrpcServer{
		QueueClient: qClient,
	})

	fmt.Println("gRPC Server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
