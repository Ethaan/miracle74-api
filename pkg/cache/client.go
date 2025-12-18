package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
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

	clientOpt, err := parseConnectionURL(cacheURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cache URL: %w", err)
	}

	client, err := valkey.NewClient(clientOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to create Valkey client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to cache: %w", err)
	}

	return &Client{
		client: client,
		ttl:    ttl,
	}, nil
}

// parseConnectionURL parses the cache URL and returns appropriate client options
// Supports formats:
// - Simple: "localhost:6379"
// - Redis URL: "redis://default:password@host:6379"
// - Redis TLS: "rediss://user:pass@host:6379"
func parseConnectionURL(cacheURL string) (valkey.ClientOption, error) {
	if strings.HasPrefix(cacheURL, "redis://") || strings.HasPrefix(cacheURL, "rediss://") {
		return parseRedisURL(cacheURL)
	}

	return valkey.ClientOption{
		InitAddress: []string{cacheURL},
	}, nil
}

// parseRedisURL parses a Redis URL with credentials
// Format: redis://[username][:password]@host:port
func parseRedisURL(redisURL string) (valkey.ClientOption, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return valkey.ClientOption{}, fmt.Errorf("invalid Redis URL: %w", err)
	}

	opt := valkey.ClientOption{
		InitAddress:  []string{u.Host},
		DisableCache: true, // Upstash doesn't support CLIENT TRACKING
	}

	if u.Scheme == "rediss" {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if u.User != nil {
		if password, ok := u.User.Password(); ok {
			opt.Password = password
		}
		if username := u.User.Username(); username != "" && username != "default" {
			opt.Username = username
		}
	}

	return opt, nil
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
