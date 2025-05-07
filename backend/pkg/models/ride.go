package models

import (
	"time"
)

type RideStatus string

const (
	RideStatusPending   RideStatus = "pending"
	RideStatusAccepted  RideStatus = "accepted"
	RideStatusStarted   RideStatus = "started"
	RideStatusCompleted RideStatus = "completed"
	RideStatusCancelled RideStatus = "cancelled"
)

type RideType string

const (
	RideTypeShared   RideType = "shared"    // Student ride sharing
	RideTypeOnDemand RideType = "on_demand" // On-demand ride like Uber/Lyft
)

type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodWallet PaymentMethod = "wallet"
)

type Ride struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	RideType       RideType      `json:"ride_type" gorm:"not null"`
	RiderID        uint          `json:"rider_id"`
	DriverID       *uint         `json:"driver_id"`
	PickupLat      float64       `json:"pickup_lat"`
	PickupLng      float64       `json:"pickup_lng"`
	DropoffLat     float64       `json:"dropoff_lat"`
	DropoffLng     float64       `json:"dropoff_lng"`
	PickupAddress  string        `json:"pickup_address"`
	DropoffAddress string        `json:"dropoff_address"`
	Status         RideStatus    `json:"status"`
	Price          float64       `json:"price"`
	Distance       float64       `json:"distance"`        // in kilometers
	Duration       int           `json:"duration"`        // in minutes
	SeatsAvailable int           `json:"seats_available"` // For shared rides
	SeatsBooked    int           `json:"seats_booked"`    // For shared rides
	DepartureTime  time.Time     `json:"departure_time"`  // For shared rides
	PaymentMethod  PaymentMethod `json:"payment_method"`
	StartedAt      *time.Time    `json:"started_at"`
	CompletedAt    *time.Time    `json:"completed_at"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`

	// Relationships
	Rider      User            `json:"rider" gorm:"foreignKey:RiderID"`
	Driver     *User           `json:"driver,omitempty" gorm:"foreignKey:DriverID"`
	Passengers []RidePassenger `json:"passengers,omitempty" gorm:"foreignKey:RideID"`
}

// RidePassenger represents a passenger in a shared ride
type RidePassenger struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	RideID    uint       `json:"ride_id"`
	UserID    uint       `json:"user_id"`
	Seats     int        `json:"seats"`
	Status    RideStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
	Ride Ride `json:"ride" gorm:"foreignKey:RideID"`
}
