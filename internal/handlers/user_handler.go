package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lilokie/otp-auth/internal/models"
	"github.com/lilokie/otp-auth/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUser handles getting a user by ID
// @Summary Get user by ID
// @Description Get a user's details by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse "User details"
// @Failure 400 {object} models.ErrorResponse "Invalid user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// Parse user ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user by ID
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return user
	response := models.UserResponse{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}
	c.JSON(http.StatusOK, response)
}

// ListUsers handles listing users with pagination and search
// @Summary List users
// @Description List users with pagination and optional search
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param search query string false "Search term for phone number"
// @Success 200 {object} models.UsersListResponse "List of users"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params.Page = 1
		params.PageSize = 10
	}

	// Set defaults if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// Get users
	users, totalCount, err := h.userService.ListUsers(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing users"})
		return
	}

	// Map to response type
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = models.UserResponse{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   user.CreatedAt,
		}
	}

	// Return response
	response := models.UsersListResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}
	c.JSON(http.StatusOK, response)
}
