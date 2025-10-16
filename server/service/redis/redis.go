package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client  *redis.Client
	enabled bool
}

func NewRedisService() (*RedisService, error) {
	// Check if Redis is enabled via environment variable
	enabled := strings.ToLower(os.Getenv("REDIS_ENABLED")) == "true"

	if !enabled {
		return &RedisService{
			client:  nil,
			enabled: false,
		}, nil
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", addr, err)
	}

	return &RedisService{
		client:  client,
		enabled: true,
	}, nil
}

func (s *RedisService) IsEnabled() bool {
	return s != nil && s.enabled
}

func (s *RedisService) Close() error {
	if !s.IsEnabled() {
		return nil
	}
	return s.client.Close()
}

// GetShortcutListCache retrieves cached shortcut list for a user
func (s *RedisService) GetShortcutListCache(ctx context.Context, userID int32) (string, error) {
	if !s.IsEnabled() {
		return "", fmt.Errorf("Redis is not enabled")
	}
	key := fmt.Sprintf("shortcuts:user:%d", userID)
	return s.client.Get(ctx, key).Result()
}

// SetShortcutListCache sets cached shortcut list for a user
func (s *RedisService) SetShortcutListCache(ctx context.Context, userID int32, data string, expiration time.Duration) error {
	if !s.IsEnabled() {
		return nil // Silently skip if Redis is disabled
	}
	key := fmt.Sprintf("shortcuts:user:%d", userID)
	return s.client.Set(ctx, key, data, expiration).Err()
}

// InvalidateShortcutListCache removes cached shortcut list for a user
func (s *RedisService) InvalidateShortcutListCache(ctx context.Context, userID int32) error {
	if !s.IsEnabled() {
		return nil // Silently skip if Redis is disabled
	}
	key := fmt.Sprintf("shortcuts:user:%d", userID)
	return s.client.Del(ctx, key).Err()
}

// SetJSON is a helper to store JSON-encoded data
func (s *RedisService) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if !s.IsEnabled() {
		return nil // Silently skip if Redis is disabled
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return s.client.Set(ctx, key, data, expiration).Err()
}

// GetJSON is a helper to retrieve and decode JSON data
func (s *RedisService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	if !s.IsEnabled() {
		return fmt.Errorf("Redis is not enabled")
	}
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
