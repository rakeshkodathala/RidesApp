package models

import (
	"time"
)

// Rating represents a rating for a ride
type Rating struct {
	ID         uint      `json:"id" gorm:"primaryKey;column:id"`
	RideID     uint      `json:"ride_id" gorm:"column:ride_id"`
	FromUserID uint      `json:"from_user_id" gorm:"column:from_user_id"`
	ToUserID   uint      `json:"to_user_id" gorm:"column:to_user_id"`
	Rating     int       `json:"rating" gorm:"column:rating"`
	Comment    string    `json:"comment" gorm:"column:comment"`
	UserID     *uint     `json:"user_id" gorm:"column:user_id"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
	FromUser   User      `json:"from_user" gorm:"foreignKey:FromUserID"`
	ToUser     User      `json:"to_user" gorm:"foreignKey:ToUserID"`
	Ride       Ride      `json:"ride" gorm:"foreignKey:RideID"`
}
