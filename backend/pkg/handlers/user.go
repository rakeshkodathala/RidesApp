package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rakeshkumar/ridesapp/pkg/models"
	"github.com/rakeshkumar/ridesapp/pkg/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userRepo: repository.NewUserRepository(),
	}
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Hash the password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user in database
	if err := h.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return created user (excluding password)
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// GetUser handles retrieving a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't send password in response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles user updates
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get existing user
	existingUser, err := h.userRepo.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Decode update data
	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update fields
	existingUser.FirstName = updateData.FirstName
	existingUser.LastName = updateData.LastName
	existingUser.Phone = updateData.Phone
	existingUser.ProfilePicture = updateData.ProfilePicture

	// If password is provided, update it
	if updateData.Password != "" {
		existingUser.Password = updateData.Password
		if err := existingUser.HashPassword(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
	}

	// Save updates
	if err := h.userRepo.UpdateUser(existingUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Return updated user (excluding password)
	existingUser.Password = ""
	c.JSON(http.StatusOK, existingUser)
}

// DeleteUser handles user deletion
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userRepo.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDrivers handles retrieving all active drivers
func (h *UserHandler) GetDrivers(c *gin.Context) {
	drivers, err := h.userRepo.GetDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve drivers"})
		return
	}

	// Don't send passwords in response
	for i := range drivers {
		drivers[i].Password = ""
	}

	c.JSON(http.StatusOK, drivers)
}

// GetCurrentUser handles retrieving the current user's profile
func GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user from database
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUserRequest represents the request body for updating a user
type UpdateCurrentUserRequest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`

	// For drivers
	LicenseNumber string `json:"license_number,omitempty"`
	VehicleModel  string `json:"vehicle_model,omitempty"`
	VehicleColor  string `json:"vehicle_color,omitempty"`
	VehiclePlate  string `json:"vehicle_plate,omitempty"`
}

// UpdateCurrentUser handles updating the current user's profile
func UpdateCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get request body
	var req UpdateCurrentUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.ProfilePicture != "" {
		user.ProfilePicture = req.ProfilePicture
	}

	// Update driver-specific fields if user is a driver
	if user.Role == models.RoleDriver {
		if req.LicenseNumber != "" {
			user.LicenseNumber = req.LicenseNumber
		}
		if req.VehicleModel != "" {
			user.VehicleModel = req.VehicleModel
		}
		if req.VehicleColor != "" {
			user.VehicleColor = req.VehicleColor
		}
		if req.VehiclePlate != "" {
			user.VehiclePlate = req.VehiclePlate
		}
	}

	// Save user to database
	if err := userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}
