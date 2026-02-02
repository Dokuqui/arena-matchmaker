package logic

import (
	"arena-matchmaker/pkg/allocator"
	"arena-matchmaker/pkg/metrics"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const PlayersPerMatch = 4

type Matchmaker struct {
	rdb        *redis.Client
	interval   time.Duration
	maxMMRDiff int
	allocator  allocator.Allocator
}

func NewMatchmaker(rdb *redis.Client, interval time.Duration, maxDiff int, alloc allocator.Allocator) *Matchmaker {
	return &Matchmaker{
		rdb:        rdb,
		interval:   interval,
		maxMMRDiff: maxDiff,
		allocator:  alloc,
	}
}

func (m *Matchmaker) Run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	fmt.Println("🔎 Matchmaker Loop Started (Party Mode: 2v2)...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("🛑 Matchmaker Loop shutting down...")
			return
		case <-ticker.C:
			m.findMatch(ctx)
		}
	}
}

func (m *Matchmaker) findMatch(ctx context.Context) {
	count, _ := m.rdb.ZCard(ctx, "matchmaker_queue").Result()
	metrics.PlayersInQueue.Set(float64(count))

	tickets, err := m.rdb.ZRangeWithScores(ctx, "matchmaker_queue", 0, 9).Result()
	if err != nil || len(tickets) < 1 {
		return
	}

	var matchTickets []string
	currentSize := 0

	for _, t := range tickets {
		ticketID := t.Member.(string)

		size, _ := m.rdb.HGet(ctx, "ticket_data:"+ticketID, "size").Int()

		if currentSize+size <= PlayersPerMatch {
			matchTickets = append(matchTickets, ticketID)
			currentSize += size
		}

		if currentSize == PlayersPerMatch {
			break
		}
	}

	if currentSize != PlayersPerMatch {
		return
	}

	pipe := m.rdb.Pipeline()
	for _, tID := range matchTickets {
		pipe.ZRem(ctx, "matchmaker_queue", tID)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("Failed to dequeue: %v", err)
		return
	}

	matchID := fmt.Sprintf("Match_%d", time.Now().Unix())
	serverIP, err := m.allocator.Allocate(ctx, matchID)
	if err != nil {
		log.Printf("Alloc failed: %v", err)
		return
	}

	metrics.MatchesMade.Inc()
	fmt.Printf("✅ TEAM MATCH! Tickets: %v\n   -> ID: %s\n   -> Server: %s\n", matchTickets, matchID, serverIP)
}
