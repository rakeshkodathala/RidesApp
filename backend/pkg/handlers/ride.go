package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakeshkumar/ridesapp/pkg/models"
	"github.com/rakeshkumar/ridesapp/pkg/repository"
)

// CreateRideRequest represents the request body for creating a ride
type CreateRideRequest struct {
	RideType       string    `json:"ride_type" binding:"required,oneof=shared on_demand"`
	PickupLat      float64   `json:"pickup_lat" binding:"required"`
	PickupLng      float64   `json:"pickup_lng" binding:"required"`
	DropoffLat     float64   `json:"dropoff_lat" binding:"required"`
	DropoffLng     float64   `json:"dropoff_lng" binding:"required"`
	PickupAddress  string    `json:"pickup_address" binding:"required"`
	DropoffAddress string    `json:"dropoff_address" binding:"required"`
	Price          float64   `json:"price" binding:"required"`
	Distance       float64   `json:"distance" binding:"required"`
	Duration       int       `json:"duration" binding:"required"`
	SeatsAvailable int       `json:"seats_available" binding:"required_if=RideType shared"`
	DepartureTime  time.Time `json:"departure_time" binding:"required_if=RideType shared"`
	PaymentMethod  string    `json:"payment_method" binding:"required,oneof=cash card wallet"`
}

// CreateRide handles the creation of a new ride
func CreateRide(c *gin.Context) {
	var req CreateRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create ride
	ride := &models.Ride{
		RideType:       models.RideType(req.RideType),
		RiderID:        userID.(uint),
		PickupLat:      req.PickupLat,
		PickupLng:      req.PickupLng,
		DropoffLat:     req.DropoffLat,
		DropoffLng:     req.DropoffLng,
		PickupAddress:  req.PickupAddress,
		DropoffAddress: req.DropoffAddress,
		Status:         models.RideStatusPending,
		Price:          req.Price,
		Distance:       req.Distance,
		Duration:       req.Duration,
		PaymentMethod:  models.PaymentMethod(req.PaymentMethod),
	}

	// Set shared ride specific fields
	if req.RideType == string(models.RideTypeShared) {
		ride.SeatsAvailable = req.SeatsAvailable
		ride.SeatsBooked = 0
		ride.DepartureTime = req.DepartureTime
	}

	// Save ride to database
	rideRepo := repository.NewRideRepository()
	if err := rideRepo.CreateRide(ride); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ride"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ride created successfully",
		"ride":    ride,
	})
}

// GetRideByID handles retrieving a ride by ID
func GetRideByID(c *gin.Context) {
	// Get ride ID from path
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// Get ride from database
	rideRepo := repository.NewRideRepository()
	ride, err := rideRepo.GetRideByID(uint(rideID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found"})
		return
	}

	c.JSON(http.StatusOK, ride)
}

// GetMyRides handles retrieving all rides for the authenticated user
func GetMyRides(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user role from context (set by auth middleware)
	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	// Get rides from database
	rideRepo := repository.NewRideRepository()
	var rides []models.Ride
	var err error

	if userRole == string(models.RoleRider) {
		rides, err = rideRepo.GetRidesByRiderID(userID.(uint))
	} else if userRole == string(models.RoleDriver) {
		rides, err = rideRepo.GetRidesByDriverID(userID.(uint))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rides"})
		return
	}

	c.JSON(http.StatusOK, rides)
}

// GetAvailableSharedRides handles retrieving all available shared rides
func GetAvailableSharedRides(c *gin.Context) {
	// Get rides from database
	rideRepo := repository.NewRideRepository()
	rides, err := rideRepo.GetAvailableSharedRides()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get available rides"})
		return
	}

	c.JSON(http.StatusOK, rides)
}

// GetUpcomingSharedRides handles retrieving upcoming shared rides
func GetUpcomingSharedRides(c *gin.Context) {
	// Get rides from database
	rideRepo := repository.NewRideRepository()
	rides, err := rideRepo.GetUpcomingSharedRides()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get upcoming rides"})
		return
	}

	c.JSON(http.StatusOK, rides)
}

// UpdateRideStatus handles updating the status of a ride
func UpdateRideStatus(c *gin.Context) {
	// Get ride ID from path
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// Get status from request body
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending accepted started completed cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update ride status in database
	rideRepo := repository.NewRideRepository()
	if err := rideRepo.UpdateRideStatus(uint(rideID), models.RideStatus(req.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ride status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ride status updated successfully"})
}

// JoinRideRequest represents the request body for joining a ride
type JoinRideRequest struct {
	Seats int `json:"seats" binding:"required,min=1"`
}

// JoinRide handles a user joining a shared ride
func JoinRide(c *gin.Context) {
	// Get ride ID from path
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get request body
	var req JoinRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create passenger
	passenger := &models.RidePassenger{
		RideID: uint(rideID),
		UserID: userID.(uint),
		Seats:  req.Seats,
		Status: models.RideStatusPending,
	}

	// Add passenger to ride
	rideRepo := repository.NewRideRepository()
	if err := rideRepo.AddPassenger(passenger); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join ride: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined ride successfully"})
}

// LeaveRide handles a user leaving a shared ride
func LeaveRide(c *gin.Context) {
	// Get passenger ID from path
	passengerID, err := strconv.ParseUint(c.Param("passengerId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passenger ID"})
		return
	}

	// Remove passenger from ride
	rideRepo := repository.NewRideRepository()
	if err := rideRepo.RemovePassenger(uint(passengerID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave ride"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left ride successfully"})
}

// GetRidePassengers handles retrieving all passengers for a ride
func GetRidePassengers(c *gin.Context) {
	// Get ride ID from path
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// Get passengers from database
	rideRepo := repository.NewRideRepository()
	passengers, err := rideRepo.GetPassengersByRideID(uint(rideID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get passengers"})
		return
	}

	c.JSON(http.StatusOK, passengers)
}

// RateRideRequest represents the request body for rating a ride
type RateRideRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

// RateRide handles rating a ride
func RateRide(c *gin.Context) {
	// Get ride ID from path
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get request body
	var req RateRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get ride to find the rider
	rideRepo := repository.NewRideRepository()
	ride, err := rideRepo.GetRideByID(uint(rideID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found"})
		return
	}

	// Create rating
	rating := &models.Rating{
		RideID:     uint(rideID),
		FromUserID: userID.(uint),
		ToUserID:   ride.RiderID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	// Save rating to database
	if err := rideRepo.AddRating(rating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating added successfully"})
}

// GetRideRatings handles retrieving all ratings for a ride
func GetRideRatings(c *gin.Context) {
	// Get ride ID from path
	rideID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// Get ratings from database
	rideRepo := repository.NewRideRepository()
	ratings, err := rideRepo.GetRatingsByRideID(uint(rideID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ratings"})
		return
	}

	c.JSON(http.StatusOK, ratings)
}
