package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisOTPRepository implements OTPRepository using Redis
type RedisOTPRepository struct {
	client *redis.Client
}

const (
	otpKeyPrefix       = "otp:"
	rateLimitKeyPrefix = "rate_limit:"
)

// NewRedisOTPRepository creates a new Redis OTP repository
func NewRedisOTPRepository(client *redis.Client) *RedisOTPRepository {
	return &RedisOTPRepository{client: client}
}

// StoreOTP stores an OTP with expiration
func (r *RedisOTPRepository) StoreOTP(ctx context.Context, phoneNumber, otp string, expiration time.Duration) error {
	key := otpKeyPrefix + phoneNumber
	err := r.client.Set(ctx, key, otp, expiration).Err()
	if err != nil {
		return fmt.Errorf("error storing OTP: %w", err)
	}
	return nil
}

// GetOTP retrieves an OTP for a phone number
func (r *RedisOTPRepository) GetOTP(ctx context.Context, phoneNumber string) (string, error) {
	key := otpKeyPrefix + phoneNumber
	otp, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("OTP not found or expired")
		}
		return "", fmt.Errorf("error retrieving OTP: %w", err)
	}
	return otp, nil
}

// DeleteOTP deletes an OTP for a phone number
func (r *RedisOTPRepository) DeleteOTP(ctx context.Context, phoneNumber string) error {
	key := otpKeyPrefix + phoneNumber
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting OTP: %w", err)
	}
	return nil
}

// CheckRateLimit checks if the rate limit for a phone number has been exceeded
func (r *RedisOTPRepository) CheckRateLimit(ctx context.Context, phoneNumber string, limit int, window time.Duration) (bool, error) {
	key := rateLimitKeyPrefix + phoneNumber
	count, err := r.client.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, fmt.Errorf("error checking rate limit: %w", err)
	}
	return count >= limit, nil
}

// IncrementRateLimit increments the rate limit counter for a phone number
func (r *RedisOTPRepository) IncrementRateLimit(ctx context.Context, phoneNumber string, window time.Duration) error {
	key := rateLimitKeyPrefix + phoneNumber

	// Check if key exists
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error checking if rate limit key exists: %w", err)
	}

	// If key doesn't exist, set it with expiration
	if exists == 0 {
		err = r.client.Set(ctx, key, 1, window).Err()
		if err != nil {
			return fmt.Errorf("error setting rate limit: %w", err)
		}
		return nil
	}

	// Otherwise, increment it
	_, err = r.client.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error incrementing rate limit: %w", err)
	}

	return nil
}
