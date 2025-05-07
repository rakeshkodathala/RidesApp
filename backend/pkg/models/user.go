package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleRider  UserRole = "rider"
	RoleDriver UserRole = "driver"
)

type User struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Email          string    `json:"email" gorm:"unique;not null"`
	Password       string    `json:"-" gorm:"not null"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	Role           UserRole  `json:"role"`
	ProfilePicture string    `json:"profile_picture"`                  // URL to profile picture
	Rating         float64   `json:"rating" gorm:"default:5.0"`        // User rating
	IsVerified     bool      `json:"is_verified" gorm:"default:false"` // Whether the user is verified
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// For drivers
	LicenseNumber string `json:"license_number,omitempty"`
	VehicleModel  string `json:"vehicle_model,omitempty"`
	VehicleColor  string `json:"vehicle_color,omitempty"`
	VehiclePlate  string `json:"vehicle_plate,omitempty"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
