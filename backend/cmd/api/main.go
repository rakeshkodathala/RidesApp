package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rakeshkumar/ridesapp/pkg/config"
	"github.com/rakeshkumar/ridesapp/pkg/database"
	"github.com/rakeshkumar/ridesapp/pkg/handlers"
	"github.com/rakeshkumar/ridesapp/pkg/middleware"
	"github.com/rakeshkumar/ridesapp/pkg/models"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Error loading config, using defaults: %v", err)
		cfg = &config.Config{
			DBHost:     "localhost",
			DBPort:     "5432",
			DBUser:     "postgres",
			DBPassword: "postgres",
			DBName:     "ridesapp",
			ServerPort: "8080",
		}
	}

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "OK",
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", handlers.GetCurrentUser)
				users.PUT("/me", handlers.UpdateCurrentUser)
			}

			// Ride routes
			rides := protected.Group("/rides")
			{
				// Create a new ride
				rides.POST("", handlers.CreateRide)

				// Get a specific ride
				rides.GET("/:id", handlers.GetRideByID)

				// Get my rides (as rider or driver)
				rides.GET("/my", handlers.GetMyRides)

				// Get available shared rides
				rides.GET("/shared/available", handlers.GetAvailableSharedRides)

				// Get upcoming shared rides
				rides.GET("/shared/upcoming", handlers.GetUpcomingSharedRides)

				// Update ride status
				rides.PUT("/:id/status", handlers.UpdateRideStatus)

				// Join a shared ride
				rides.POST("/:id/join", handlers.JoinRide)

				// Leave a shared ride
				rides.DELETE("/:id/passengers/:passengerId", handlers.LeaveRide)

				// Get passengers for a ride
				rides.GET("/:id/passengers", handlers.GetRidePassengers)
			}

			// Driver routes
			drivers := protected.Group("/drivers")
			drivers.Use(middleware.RoleMiddleware(models.RoleDriver))
			{
				// Driver-specific endpoints can be added here
			}
		}
	}

	// Get port from environment variable or use default
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
