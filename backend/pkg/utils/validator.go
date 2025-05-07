package utils

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// IsValidEmail checks if the email is valid
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidPhone checks if the phone number is valid
func IsValidPhone(phone string) bool {
	// Basic phone validation - can be enhanced based on requirements
	return len(phone) >= 10
}

// IsValidPassword checks if the password meets requirements
func IsValidPassword(password string) bool {
	// Password should be at least 8 characters
	return len(password) >= 8
}
