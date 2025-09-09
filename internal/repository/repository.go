package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lilokie/otp-auth/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, phoneNumber string) (*models.User, error)

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)

	// FindByPhoneNumber finds a user by phone number
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)

	// List returns a list of users with pagination and search
	List(ctx context.Context, params models.PaginationParams) ([]models.User, int64, error)

	// Update updates a user
	Update(ctx context.Context, user *models.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
}

// OTPRepository defines the interface for OTP operations
type OTPRepository interface {
	// StoreOTP stores an OTP with expiration
	StoreOTP(ctx context.Context, phoneNumber, otp string, expiration time.Duration) error

	// GetOTP retrieves an OTP for a phone number
	GetOTP(ctx context.Context, phoneNumber string) (string, error)

	// DeleteOTP deletes an OTP for a phone number
	DeleteOTP(ctx context.Context, phoneNumber string) error

	// CheckRateLimit checks if the rate limit for a phone number has been exceeded
	CheckRateLimit(ctx context.Context, phoneNumber string, limit int, window time.Duration) (bool, error)

	// IncrementRateLimit increments the rate limit counter for a phone number
	IncrementRateLimit(ctx context.Context, phoneNumber string, window time.Duration) error
}
