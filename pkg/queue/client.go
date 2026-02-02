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

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	return &Client{RDB: rdb}, nil
}

func (c *Client) AddToQueue(ctx context.Context, ticketID string, playerIDs []string, avgMMR int) error {
	pipe := c.RDB.Pipeline()

	membersStr := fmt.Sprintf("%v", playerIDs)

	key := "ticket_data:" + ticketID
	pipe.HSet(ctx, key, "size", len(playerIDs))
	pipe.HSet(ctx, key, "members", membersStr)

	pipe.ZAdd(ctx, "matchmaker_queue", redis.Z{
		Score:  float64(avgMMR),
		Member: ticketID,
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to queue ticket: %w", err)
	}
	return nil
}

func (c *Client) GetTicketSize(ctx context.Context, ticketID string) (int, error) {
	val, err := c.RDB.HGet(ctx, "ticket_data:"+ticketID, "size").Int()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) GetPlayerMMR(ctx context.Context, playerID string) (int, error) {
	val, err := c.RDB.HGet(ctx, "player_ratings", playerID).Int()
	if err == redis.Nil {
		return 1000, nil
	}
	return val, err
}

func (c *Client) SetPlayerMMR(ctx context.Context, playerID string, mmr int) error {
	return c.RDB.HSet(ctx, "player_ratings", playerID, mmr).Err()
}
