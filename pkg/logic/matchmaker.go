package logic

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	TickRate        = 1 * time.Second
	PlayersPerMatch = 2
	MaxMMRDiff      = 200
)

type Matchmaker struct {
	rdb *redis.Client
}

func NewMatchmaker(rdb *redis.Client) *Matchmaker {
	return &Matchmaker{rdb: rdb}
}
func (m *Matchmaker) Run(ctx context.Context) {
	ticker := time.NewTicker(TickRate)
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

	if mmrDiff > MaxMMRDiff {
		// Optimization: In a real system, we might expand the search here.
		// For now, we just wait for someone better to join.
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
	fmt.Printf("MATCH CREATED! [%s] vs [%s] (Diff: %.0f) -> ID: %s\n",
		p1.Member, p2.Member, mmrDiff, matchID)
}
