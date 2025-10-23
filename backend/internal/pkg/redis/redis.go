package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func New() *Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Client{rdb}
}

// CacheSet sets key with value (JSON) and TTL (1 hour)
func (c *Client) CacheSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, data, ttl).Err()
}

// CacheGet gets value by key (JSON unmarshal to v)
func (c *Client) CacheGet(ctx context.Context, key string, v interface{}) error {
	data, err := c.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("cache miss: %w", err)
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// CacheDel deletes key
func (c *Client) CacheDel(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}
