package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/config"
	"arena-matchmaker/pkg/queue"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Arena Frontend (API)...")

	cfg := config.Load()

	fmt.Printf("🔌 Connecting to Redis at: %s\n", cfg.RedisAddr)
	qClient, err := queue.NewClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Could not initialize Redis: %v", err)
	}
	defer qClient.RDB.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("📊 Metrics listening on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Printf("Metrics server failed: %v", err)
		}
	}()

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
	fmt.Println("\n⚠️  Shutting down Frontend...")
	grpcServer.GracefulStop()
	fmt.Println("✅ Frontend stopped.")
}
