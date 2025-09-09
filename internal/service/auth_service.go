package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lilokie/otp-auth/config"
	"github.com/lilokie/otp-auth/internal/models"
	"github.com/lilokie/otp-auth/internal/repository"
)

// AuthService handles authentication-related business logic
type AuthService struct {
	userRepo repository.UserRepository
	otpRepo  repository.OTPRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	otpRepo repository.OTPRepository,
	config *config.Config,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		otpRepo:  otpRepo,
		config:   config,
	}
}

// GenerateOTP generates a one-time password for a phone number
func (s *AuthService) GenerateOTP(ctx context.Context, phoneNumber string) (string, error) {
	// Check rate limit
	exceeded, err := s.otpRepo.CheckRateLimit(ctx, phoneNumber, s.config.OTP.RateLimit.Count, s.config.GetRateLimitDuration())
	if err != nil {
		return "", fmt.Errorf("error checking rate limit: %w", err)
	}
	if exceeded {
		return "", fmt.Errorf("rate limit exceeded")
	}

	// Generate OTP
	otp := s.generateRandomOTP(s.config.OTP.Length)

	// Store OTP in Redis
	err = s.otpRepo.StoreOTP(ctx, phoneNumber, otp, s.config.GetOTPExpiration())
	if err != nil {
		return "", fmt.Errorf("error storing OTP: %w", err)
	}

	// Increment rate limit
	err = s.otpRepo.IncrementRateLimit(ctx, phoneNumber, s.config.GetRateLimitDuration())
	if err != nil {
		return "", fmt.Errorf("error incrementing rate limit: %w", err)
	}

	return otp, nil
}

// VerifyOTP verifies an OTP and returns a JWT token if valid
func (s *AuthService) VerifyOTP(ctx context.Context, phoneNumber, otp string) (string, *models.User, error) {
	// Get stored OTP
	storedOTP, err := s.otpRepo.GetOTP(ctx, phoneNumber)
	if err != nil {
		return "", nil, fmt.Errorf("error retrieving OTP: %w", err)
	}

	// Verify OTP
	if storedOTP != otp {
		return "", nil, fmt.Errorf("invalid OTP")
	}

	// Delete OTP to prevent reuse
	err = s.otpRepo.DeleteOTP(ctx, phoneNumber)
	if err != nil {
		return "", nil, fmt.Errorf("error deleting OTP: %w", err)
	}

	// Find user by phone number or create if not exists
	user, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		// User not found, create new user
		user, err = s.userRepo.Create(ctx, phoneNumber)
		if err != nil {
			return "", nil, fmt.Errorf("error creating user: %w", err)
		}
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("error generating JWT: %w", err)
	}

	return token, user, nil
}

// generateRandomOTP generates a random numeric OTP of the specified length
func (s *AuthService) generateRandomOTP(length int) string {
	// Use a proper random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random number with specified length
	min := int64(1) * int64(powInt(10, length-1))
	max := int64(9) * int64(powInt(10, length-1))

	otpNum := min + r.Int63n(max-min+1)
	return strconv.FormatInt(otpNum, 10)
}

// powInt calculates x^y
func powInt(x, y int) int {
	result := 1
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}

// generateJWT generates a JWT token for a user
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	// Create the JWT claims, which includes the user ID and expiry time
	expirationTime := time.Now().Add(time.Duration(s.config.JWT.ExpirationHours) * time.Hour)

	claims := jwt.MapClaims{
		"user_id":      user.ID.String(),
		"phone_number": user.PhoneNumber,
		"exp":          expirationTime.Unix(),
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
