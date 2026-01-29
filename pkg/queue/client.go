package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
}

func NewClient(addr string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{RDB: rdb}, nil
}

func (c *Client) AddToQueue(ctx context.Context, playerID string, mmr int32) error {
	err := c.RDB.ZAdd(ctx, "matchmaker_queue", redis.Z{
		Score:  float64(mmr),
		Member: playerID,
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to add player to queue: %w", err)
	}
	return nil
}

func (c *Client) GetPlayerMMR(ctx context.Context, playerID string) (int, error) {
	val, err := c.RDB.HGet(ctx, "player_ratings", playerID).Int()
	if err == redis.Nil {
		return 1000, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get mmr: %w", err)
	}
	return val, nil
}

func (c *Client) SetPlayerMMR(ctx context.Context, playerID string, mmr int) error {
	err := c.RDB.HSet(ctx, "player_ratings", playerID, mmr).Err()
	if err != nil {
		return fmt.Errorf("failed to set mmr: %w", err)
	}
	return nil
}
