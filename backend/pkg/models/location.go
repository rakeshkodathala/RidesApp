package models

import (
	"time"
)

type Location struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Heading   float64   `json:"heading"` // in degrees
	Speed     float64   `json:"speed"`   // in km/h
	Accuracy  float64   `json:"accuracy"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationship
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type LocationUpdate struct {
	UserID    uint    `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   float64 `json:"heading"`
	Speed     float64 `json:"speed"`
	Accuracy  float64 `json:"accuracy"`
}
