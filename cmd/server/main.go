package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/queue"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Arena Matchmaker...")

	redisClient, err := queue.NewClient("localhost:6379")
	if err != nil {
		log.Fatalf("Could not initialize Redis: %v", err)
	}
	defer redisClient.Close()

	fmt.Println("Successfully connected to Redis!")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	matchmaker.RegisterMatchmakerServiceServer(grpcServer, &GrpcServer{})

	fmt.Println("gRPC Server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
