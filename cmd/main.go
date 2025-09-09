package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/lilokie/otp-auth/config"
	_ "github.com/lilokie/otp-auth/docs" // Import swagger docs
	"github.com/lilokie/otp-auth/internal/handlers"
	"github.com/lilokie/otp-auth/internal/middleware"
	"github.com/lilokie/otp-auth/internal/repository"
	"github.com/lilokie/otp-auth/internal/service"
	"github.com/lilokie/otp-auth/internal/utils"
)

// @title OTP Authentication API
// @version 1.0
// @description A RESTful API for OTP-based authentication
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
func main() {
	cfg := config.LoadConfig()

	// Setup database
	db, err := utils.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Setup Redis
	redisClient, err := utils.SetupRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to setup Redis: %v", err)
	}

	// Create repositories
	userRepo := repository.NewPostgresUserRepository(db)
	otpRepo := repository.NewRedisOTPRepository(redisClient)

	// Create services
	authService := service.NewAuthService(userRepo, otpRepo, cfg)
	userService := service.NewUserService(userRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Create middleware
	jwtMiddleware := middleware.NewJWTAuthMiddleware(cfg)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(redisClient)

	// Setup Gin router
	router := gin.Default()
	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Routes
	v1 := router.Group("/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/request-otp",
				rateLimitMiddleware.OTPRateLimit(cfg.OTP.RateLimit.Count, cfg.GetRateLimitDuration()),
				authHandler.RequestOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(jwtMiddleware.AuthRequired())
		{
			users.GET("/:id", userHandler.GetUser)
			users.GET("", userHandler.ListUsers)
		}
	}

	// Load HTML template
	tmpl, err := template.ParseFiles(filepath.Join("internal", "templates", "index.html"))
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Root route - HTML welcome page with link to Swagger UI
	rootHandler := func(c *gin.Context) {
		baseURL := fmt.Sprintf("http://%s:%s", c.Request.Host, cfg.Service.HTTP.Port)
		if err := tmpl.Execute(c.Writer, gin.H{"BaseURL": baseURL}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	}
	router.GET("/", rootHandler)
	router.HEAD("/", rootHandler)

	// API info route
	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "OTP Authentication API",
			"version":     "1.0.0",
			"description": "A RESTful API for OTP-based authentication",
			"endpoints": []gin.H{
				{"path": "/v1/auth/request-otp", "method": "POST", "description": "Request OTP for a phone number"},
				{"path": "/v1/auth/verify-otp", "method": "POST", "description": "Verify OTP for a phone number"},
				{"path": "/v1/users/:id", "method": "GET", "description": "Get user by ID"},
				{"path": "/v1/users", "method": "GET", "description": "List users with pagination and search"},
			},
			"docs_url": "/swagger/index.html",
		})
	})

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Service.HTTP.Port),
		Handler: router,
	}

	// Run server in a goroutine so it doesn't block
	go func() {
		log.Printf("Server starting on port %s", cfg.Service.HTTP.Port)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown using config
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetGracefulShutdownDuration())
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database and Redis connections
	log.Println("Closing database connection...")
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Closing Redis connection...")
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Server exited properly")
}
