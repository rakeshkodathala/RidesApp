package main

import (
	"log"

	"github.com/rakeshkumar/ridesapp/pkg/config"
	"github.com/rakeshkumar/ridesapp/pkg/database"
	"github.com/rakeshkumar/ridesapp/pkg/models"
	"github.com/rakeshkumar/ridesapp/pkg/repository"
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

	// Create user repository
	userRepo := repository.NewUserRepository()

	// Create a test user
	testUser := &models.User{
		Email:     "test@example.com",
		Password:  "test123",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "+1234567890",
		Role:      models.RoleRider,
	}

	// Hash the password
	err = testUser.HashPassword()
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create the user in the database
	err = userRepo.CreateUser(testUser)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	log.Printf("Successfully created test user with ID: %d", testUser.ID)

	// Verify the user was created by retrieving it
	retrievedUser, err := userRepo.GetUserByEmail(testUser.Email)
	if err != nil {
		log.Fatalf("Failed to retrieve user: %v", err)
	}

	log.Printf("Retrieved user: %+v", retrievedUser)
}
