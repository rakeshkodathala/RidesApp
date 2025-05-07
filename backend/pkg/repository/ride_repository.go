package repository

import (
	"errors"
	"time"

	"github.com/rakeshkumar/ridesapp/pkg/database"
	"github.com/rakeshkumar/ridesapp/pkg/models"
	"gorm.io/gorm"
)

type RideRepository struct {
	db *gorm.DB
}

func NewRideRepository() *RideRepository {
	return &RideRepository{
		db: database.GetDB(),
	}
}

// CreateRide creates a new ride in the database
func (r *RideRepository) CreateRide(ride *models.Ride) error {
	return r.db.Create(ride).Error
}

// GetRideByID retrieves a ride by ID
func (r *RideRepository) GetRideByID(id uint) (*models.Ride, error) {
	var ride models.Ride
	if err := r.db.Preload("Rider").Preload("Driver").Preload("Passengers.User").First(&ride, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ride not found")
		}
		return nil, err
	}
	return &ride, nil
}

// GetRidesByRiderID retrieves all rides for a specific rider
func (r *RideRepository) GetRidesByRiderID(riderID uint) ([]models.Ride, error) {
	var rides []models.Ride
	if err := r.db.Where("rider_id = ?", riderID).Preload("Rider").Preload("Driver").Preload("Passengers.User").Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

// GetRidesByDriverID retrieves all rides for a specific driver
func (r *RideRepository) GetRidesByDriverID(driverID uint) ([]models.Ride, error) {
	var rides []models.Ride
	if err := r.db.Where("driver_id = ?", driverID).Preload("Rider").Preload("Driver").Preload("Passengers.User").Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

// GetAvailableSharedRides retrieves all available shared rides
func (r *RideRepository) GetAvailableSharedRides() ([]models.Ride, error) {
	var rides []models.Ride
	if err := r.db.Where("ride_type = ? AND status = ? AND seats_available > seats_booked",
		models.RideTypeShared, models.RideStatusPending).
		Preload("Rider").
		Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

// GetUpcomingSharedRides retrieves upcoming shared rides
func (r *RideRepository) GetUpcomingSharedRides() ([]models.Ride, error) {
	var rides []models.Ride
	now := time.Now()
	if err := r.db.Where("ride_type = ? AND status = ? AND departure_time > ?",
		models.RideTypeShared, models.RideStatusPending, now).
		Preload("Rider").
		Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

// UpdateRide updates a ride in the database
func (r *RideRepository) UpdateRide(ride *models.Ride) error {
	return r.db.Save(ride).Error
}

// UpdateRideStatus updates the status of a ride
func (r *RideRepository) UpdateRideStatus(rideID uint, status models.RideStatus) error {
	return r.db.Model(&models.Ride{}).Where("id = ?", rideID).Update("status", status).Error
}

// AddPassenger adds a passenger to a shared ride
func (r *RideRepository) AddPassenger(passenger *models.RidePassenger) error {
	// Start a transaction
	tx := r.db.Begin()

	// Get the ride
	var ride models.Ride
	if err := tx.First(&ride, passenger.RideID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if there are enough seats available
	if ride.SeatsBooked+passenger.Seats > ride.SeatsAvailable {
		tx.Rollback()
		return errors.New("not enough seats available")
	}

	// Update the ride's booked seats
	if err := tx.Model(&ride).Update("seats_booked", ride.SeatsBooked+passenger.Seats).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add the passenger
	if err := tx.Create(passenger).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// RemovePassenger removes a passenger from a shared ride
func (r *RideRepository) RemovePassenger(passengerID uint) error {
	// Start a transaction
	tx := r.db.Begin()

	// Get the passenger
	var passenger models.RidePassenger
	if err := tx.First(&passenger, passengerID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Get the ride
	var ride models.Ride
	if err := tx.First(&ride, passenger.RideID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the ride's booked seats
	if err := tx.Model(&ride).Update("seats_booked", ride.SeatsBooked-passenger.Seats).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the passenger
	if err := tx.Delete(&passenger).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// GetPassengersByRideID retrieves all passengers for a specific ride
func (r *RideRepository) GetPassengersByRideID(rideID uint) ([]models.RidePassenger, error) {
	var passengers []models.RidePassenger
	if err := r.db.Where("ride_id = ?", rideID).Preload("User").Find(&passengers).Error; err != nil {
		return nil, err
	}
	return passengers, nil
}

// AddRating adds a rating for a ride
func (r *RideRepository) AddRating(rating *models.Rating) error {
	return r.db.Create(rating).Error
}

// GetRatingsByRideID retrieves all ratings for a specific ride
func (r *RideRepository) GetRatingsByRideID(rideID uint) ([]models.Rating, error) {
	var ratings []models.Rating
	result := r.db.Debug().
		Where("ride_id = ?", rideID).
		Preload("FromUser").
		Preload("ToUser").
		Preload("Ride").
		Find(&ratings)

	if result.Error != nil {
		return nil, result.Error
	}
	return ratings, nil
}
