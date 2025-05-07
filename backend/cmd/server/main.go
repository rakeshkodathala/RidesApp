package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rakeshkumar/ridesapp/pkg/config"
	"github.com/rakeshkumar/ridesapp/pkg/database"
	"github.com/rakeshkumar/ridesapp/pkg/handlers"
	"github.com/rakeshkumar/ridesapp/pkg/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	err = database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is working!",
		})
	})

	// Auth routes
	router.POST("/api/v1/auth/register", handlers.Register)
	router.POST("/api/v1/auth/login", handlers.Login)

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/me", handlers.GetCurrentUser)
		protected.PUT("/users/me", handlers.UpdateCurrentUser)

		// Ride routes - Static paths first
		protected.POST("/rides", handlers.CreateRide)
		protected.GET("/rides/my-rides", handlers.GetMyRides)
		protected.GET("/rides/shared/available", handlers.GetAvailableSharedRides)
		protected.GET("/rides/shared/upcoming", handlers.GetUpcomingSharedRides)

		// Ride routes - Dynamic paths with parameters
		protected.GET("/rides/:id", handlers.GetRideByID)
		protected.PUT("/rides/:id/status", handlers.UpdateRideStatus)
		protected.POST("/rides/:id/join", handlers.JoinRide)
		protected.GET("/rides/:id/passengers", handlers.GetRidePassengers)
		protected.DELETE("/rides/:id/passengers/:passengerId", handlers.LeaveRide)
		protected.POST("/rides/:id/rate", handlers.RateRide)
		protected.GET("/rides/:id/ratings", handlers.GetRideRatings)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
