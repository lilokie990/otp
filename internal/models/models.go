package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// OTP represents a one-time password
type OTP struct {
	PhoneNumber string    `json:"phone_number"`
	Code        string    `json:"code"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// RequestOTPRequest is the request to get an OTP
type RequestOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// RequestOTPResponse is the response to an OTP request
type RequestOTPResponse struct {
	Message string `json:"message"` // OTP is now only printed to console logs
}

// VerifyOTPRequest is the request to verify an OTP
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required,len=6,numeric"`
}

// VerifyOTPResponse is the response to an OTP verification
type VerifyOTPResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserResponse is the response containing user information
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

// UsersListResponse is the response for listing users
type UsersListResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

// PaginationParams defines pagination parameters for listing users
type PaginationParams struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Search   string `form:"search" json:"search"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// TokenClaims represents the custom JWT claims
type TokenClaims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
}
