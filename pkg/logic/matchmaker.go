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

const PlayersPerMatch = 2

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

	fmt.Println("Matchmaker Loop Started...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping Matchmaker...")
			return
		case <-ticker.C:
			m.findMatch(ctx)
		}
	}
}

func (m *Matchmaker) findMatch(ctx context.Context) {
	count, _ := m.rdb.ZCard(ctx, "matchmaker_queue").Result()
	metrics.PlayersInQueue.Set(float64(count))

	players, err := m.rdb.ZRangeWithScores(ctx, "matchmaker_queue", 0, 1).Result()
	if err != nil {
		log.Printf("Error reading queue: %v", err)
		return
	}

	if len(players) < PlayersPerMatch {
		return
	}

	p1 := players[0]
	p2 := players[1]

	mmrDiff := p1.Score - p2.Score
	if mmrDiff < 0 {
		mmrDiff = -mmrDiff
	}

	if mmrDiff > float64(m.maxMMRDiff) {
		return
	}

	pipe := m.rdb.Pipeline()
	pipe.ZRem(ctx, "matchmaker_queue", p1.Member)
	pipe.ZRem(ctx, "matchmaker_queue", p2.Member)
	_, err = pipe.Exec(ctx)

	if err != nil {
		log.Printf("Failed to remove players from queue: %v", err)
		return
	}

	matchID := fmt.Sprintf("Match_%d", time.Now().Unix())

	serverIP, err := m.allocator.Allocate(ctx, matchID)
	if err != nil {
		log.Printf("❌ Failed to allocate server: %v", err)
		// In a real system, we would re-queue the players here!
		return
	}

	metrics.MatchesMade.Inc()

	fmt.Printf("✅ MATCH READY! [%s] vs [%s] (Diff: %.0f)\n   -> ID: %s\n   -> Server: %s\n",
		p1.Member, p2.Member, mmrDiff, matchID, serverIP)
}
