package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(addr string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) AddToQueue(ctx context.Context, playerID string, mmr int32) error {
	err := c.rdb.ZAdd(ctx, "matchmaker_queue", redis.Z{
		Score:  float64(mmr),
		Member: playerID,
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to add player to queue: %w", err)
	}
	return nil
}
