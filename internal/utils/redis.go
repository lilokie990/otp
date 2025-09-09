package utils

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/lilokie/otp-auth/config"
)

// SetupRedis sets up the Redis connection
func SetupRedis(config *config.Config) (*redis.Client, error) {
	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	return client, nil
}
