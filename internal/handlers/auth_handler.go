package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lilokie/otp-auth/internal/models"
	"github.com/lilokie/otp-auth/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RequestOTP handles OTP request
// @Summary Request OTP for a phone number
// @Description Generate and send a one-time password to the provided phone number (OTP is printed to server logs)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RequestOTPRequest true "Phone number to send OTP to"
// @Success 200 {object} models.RequestOTPResponse "OTP sent successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 429 {object} models.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/request-otp [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req models.RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format")})
		return
	}

	phoneNumber := req.PhoneNumber
	// Allow any non-empty phone number for testing purposes
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number cannot be empty"})
		return
	}

	// Validate Iranian phone number format: must start with +98, 98, or 09 and be 13, 12, or 11 digits respectively
	if !(strings.HasPrefix(phoneNumber, "+98") && len(phoneNumber) == 13) &&
		!(strings.HasPrefix(phoneNumber, "98") && len(phoneNumber) == 12) &&
		!(strings.HasPrefix(phoneNumber, "09") && len(phoneNumber) == 11) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Iranian phone number format. Use +989XXXXXXXXX, 989XXXXXXXXX, or 09XXXXXXXXX"})
		return
	}

	// Generate OTP
	otp, err := h.authService.GenerateOTP(c.Request.Context(), phoneNumber)
	if err != nil {
		if err.Error() == "rate limit exceeded" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error generating OTP: %v", err)})
		return
	}

	// Print OTP to console log instead of returning it in the response
	fmt.Printf("[OTP] Phone: %s, Code: %s\n", phoneNumber, otp)

	// Return response without OTP
	response := models.RequestOTPResponse{
		Message: "OTP sent successfully. Check server logs for the code.",
	}
	c.JSON(http.StatusOK, response)
}

// VerifyOTP handles OTP verification
// @Summary Verify OTP for a phone number
// @Description Verify the OTP provided for a phone number and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Phone number and OTP to verify"
// @Success 200 {object} models.VerifyOTPResponse "OTP verified successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Invalid or expired OTP"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorMessage := "Invalid request format"
		// Provide more specific error message based on validation failure
		if err.Error() == "Key: 'VerifyOTPRequest.PhoneNumber' Error:Field validation for 'PhoneNumber' failed on the 'iranianMobile' tag" {
			errorMessage = "Invalid phone number format. Use Iranian mobile format: +989XXXXXXXXX, 09XXXXXXXXX, or 9XXXXXXXXX"
		} else if err.Error() == "Key: 'VerifyOTPRequest.OTP' Error:Field validation for 'OTP' failed on the 'len' tag" {
			errorMessage = "OTP must be exactly 6 digits"
		} else if err.Error() == "Key: 'VerifyOTPRequest.OTP' Error:Field validation for 'OTP' failed on the 'numeric' tag" {
			errorMessage = "OTP must contain only numbers"
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	phoneNumber := req.PhoneNumber
	// Allow any non-empty phone number for testing purposes
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number cannot be empty"})
		return
	}
	// Validate Iranian phone number format: must start with +98, 98, or 09 and be 13, 12, or 11 digits respectively
	if !(strings.HasPrefix(phoneNumber, "+98") && len(phoneNumber) == 13) &&
		!(strings.HasPrefix(phoneNumber, "98") && len(phoneNumber) == 12) &&
		!(strings.HasPrefix(phoneNumber, "09") && len(phoneNumber) == 11) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Iranian phone number format. Use +989XXXXXXXXX, 989XXXXXXXXX, or 09XXXXXXXXX"})
		return
	}

	// Verify OTP
	token, user, err := h.authService.VerifyOTP(c.Request.Context(), phoneNumber, req.OTP)
	if err != nil {
		if err.Error() == "invalid OTP" || err.Error() == "error retrieving OTP: OTP not found or expired" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error verifying OTP: %v", err)})
		return
	}

	// Return response
	response := models.VerifyOTPResponse{
		Token: token,
		User:  *user,
	}
	c.JSON(http.StatusOK, response)
}
