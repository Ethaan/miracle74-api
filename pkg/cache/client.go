package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

const (
	DefaultTTL = 10 * time.Minute
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Client struct {
	client valkey.Client
	ttl    time.Duration
}

func NewClient(cacheURL string, ttl time.Duration) (*Client, error) {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{cacheURL},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Valkey client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Valkey: %w", err)
	}

	return &Client{
		client: client,
		ttl:    ttl,
	}, nil
}

func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	result, err := c.client.Do(ctx, c.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal([]byte(result), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}) error {
	return c.SetWithTTL(ctx, key, value, c.ttl)
}

func (c *Client) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	cmd := c.client.B().Set().Key(key).Value(string(data)).ExSeconds(int64(ttl.Seconds())).Build()
	if err := c.client.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	cmd := c.client.B().Del().Key(key).Build()
	if err := c.client.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

func (c *Client) Close() {
	c.client.Close()
}
