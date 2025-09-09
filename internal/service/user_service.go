package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lilokie/otp-auth/internal/models"
	"github.com/lilokie/otp-auth/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	return user, nil
}

// GetUserByPhoneNumber gets a user by phone number
func (s *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	user, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("error getting user by phone number: %w", err)
	}
	return user, nil
}

// ListUsers lists users with pagination and search
func (s *UserService) ListUsers(ctx context.Context, params models.PaginationParams) ([]models.User, int64, error) {
	users, totalCount, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing users: %w", err)
	}
	return users, totalCount, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	err := s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
