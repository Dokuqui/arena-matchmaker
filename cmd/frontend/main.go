package main

import (
	"arena-matchmaker/gen/matchmaker"
	"arena-matchmaker/pkg/config"
	"arena-matchmaker/pkg/queue"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DashboardStats struct {
	EUCount int64 `json:"eu_count"`
	USCount int64 `json:"us_count"`
}

func startDashboardServer(q *queue.Client) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WS Upgrade error: %v", err)
			return
		}
		defer conn.Close()

		for {
			euCount, _ := q.RDB.ZCard(context.Background(), "queue:EU-WEST").Result()
			usCount, _ := q.RDB.ZCard(context.Background(), "queue:US-EAST").Result()

			stats := DashboardStats{
				EUCount: euCount,
				USCount: usCount,
			}

			msg, _ := json.Marshal(stats)
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}

			time.Sleep(500 * time.Millisecond)
		}
	})

	fmt.Println("💻 Dashboard running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Dashboard failed: %v", err)
	}
}

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
		http.ListenAndServe(":2112", nil)
	}()

	go startDashboardServer(qClient)

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
}
