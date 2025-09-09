package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimitMiddleware is a middleware for rate limiting
type RateLimitMiddleware struct {
	redisClient *redis.Client
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(redisClient *redis.Client) *RateLimitMiddleware {
	return &RateLimitMiddleware{redisClient: redisClient}
}

// RateLimit limits the number of requests based on IP address
func (m *RateLimitMiddleware) RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get IP address
		ip := c.ClientIP()
		key := "rate_limit:ip:" + ip

		// Check if key exists
		ctx := c.Request.Context()
		val, err := m.redisClient.Get(ctx, key).Int()
		if err != nil && !errors.Is(err, redis.Nil) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking rate limit"})
			c.Abort()
			return
		}

		// If key doesn't exist, set it
		if errors.Is(err, redis.Nil) {
			m.redisClient.Set(ctx, key, 1, window)
			c.Next()
			return
		}

		// Check if limit is exceeded
		if val >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		// Increment counter
		m.redisClient.Incr(ctx, key)

		// Continue with request
		c.Next()
	}
}

// OTPRateLimit specifically limits OTP request rate by phone number and IP address
// This provides stronger protection against OTP abuse by limiting both per-IP and per-phone number
// OTPRateLimit specifically limits OTP request rate by phone number and IP address
func (m *RateLimitMiddleware) OTPRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check IP-based rate limit (basic protection)
		ip := c.ClientIP()
		ipKey := fmt.Sprintf("rate_limit:otp:ip:%s", ip)

		// Read and preserve the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read request body"})
			c.Abort()
			return
		}

		// Restore the body so it can be read again
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Try to extract phone number from request body
		var requestBody struct {
			PhoneNumber string `json:"phone_number"`
		}

		phoneBasedLimiting := false
		phoneKey := ""

		if err := json.Unmarshal(bodyBytes, &requestBody); err == nil && requestBody.PhoneNumber != "" {
			phoneBasedLimiting = true
			phoneKey = fmt.Sprintf("rate_limit:otp:phone:%s", requestBody.PhoneNumber)
		}

		ctx := c.Request.Context()

		// Check IP-based rate limit
		ipCount, err := m.redisClient.Get(ctx, ipKey).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking rate limit"})
			c.Abort()
			return
		}

		// If IP key doesn't exist, set it
		if err == redis.Nil {
			m.redisClient.Set(ctx, ipKey, 1, window)
		} else {
			// If IP limit is exceeded
			if ipCount >= limit*2 { // IP limit is higher than phone number limit
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				c.Abort()
				return
			}
			// Increment IP counter
			m.redisClient.Incr(ctx, ipKey)
		}

		// If we can do phone-based limiting
		if phoneBasedLimiting {
			// Check phone-based rate limit
			phoneCount, err := m.redisClient.Get(ctx, phoneKey).Int()
			if err != nil && err != redis.Nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking rate limit"})
				c.Abort()
				return
			}

			// If phone key doesn't exist, set it
			if err == redis.Nil {
				m.redisClient.Set(ctx, phoneKey, 1, window)
			} else {
				// If phone limit is exceeded
				if phoneCount >= limit {
					c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many OTP requests for this phone number"})
					c.Abort()
					return
				}
				// Increment phone counter
				m.redisClient.Incr(ctx, phoneKey)
			}
		}

		// Continue with request
		c.Next()
	}
}
