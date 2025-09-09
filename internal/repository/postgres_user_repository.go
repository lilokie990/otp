package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lilokie/otp-auth/internal/models"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := `
		INSERT INTO users (id, phone_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, phone_number, created_at, updated_at
	`

	now := time.Now()
	id := uuid.New()

	user := &models.User{}
	err := r.db.QueryRowxContext(
		ctx,
		query,
		id,
		phoneNumber,
		now,
		now,
	).StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

// FindByID finds a user by ID
func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, phone_number, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, fmt.Errorf("error finding user by ID: %w", err)
	}

	return user, nil
}

// FindByPhoneNumber finds a user by phone number
func (r *PostgresUserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := `
		SELECT id, phone_number, created_at, updated_at
		FROM users
		WHERE phone_number = $1
	`

	user := &models.User{}
	err := r.db.GetContext(ctx, user, query, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("error finding user by phone number: %w", err)
	}

	return user, nil
}

// List returns a list of users with pagination and search
func (r *PostgresUserRepository) List(ctx context.Context, params models.PaginationParams) ([]models.User, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize

	// Base query
	countQuery := `SELECT COUNT(*) FROM users`
	query := `
		SELECT id, phone_number, created_at, updated_at
		FROM users
	`

	// Add search condition if provided
	var args []interface{}
	if params.Search != "" {
		whereClause := `WHERE phone_number LIKE $1`
		countQuery = countQuery + " " + whereClause
		query = query + " " + whereClause
		args = append(args, "%"+params.Search+"%")
	}

	// Add pagination
	query = query + ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) +
		` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	args = append(args, params.PageSize, offset)

	// Get total count
	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, args[:len(args)-2]...)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting users: %w", err)
	}

	// Get users
	var users []models.User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing users: %w", err)
	}

	return users, totalCount, nil
}

// Update updates a user
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET phone_number = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.PhoneNumber,
		now,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	user.UpdatedAt = now
	return nil
}

// Delete deletes a user
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
