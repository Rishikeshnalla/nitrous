package models

import "time"

// Event represents a racing event
type Event struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Location     string    `json:"location"`
	Date         time.Time `json:"date"`
	Time         string    `json:"time,omitempty"`
	IsLive       bool      `json:"isLive"`
	Category     string    `json:"category"`
	ThumbnailURL string    `json:"thumbnailUrl,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Category represents an event category
type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Icon        string `json:"icon"`
	LiveCount   int    `json:"liveCount"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Journey represents an exclusive experience
type Journey struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	Badge        string    `json:"badge"`
	SlotsLeft    int       `json:"slotsLeft"`
	Date         time.Time `json:"date"`
	Price        float64   `json:"price"`
	ThumbnailURL string    `json:"thumbnailUrl,omitempty"`
}

// MerchItem represents a merchandise item
type MerchItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Icon     string  `json:"icon"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// User represents a platform user
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never send password hash to frontend
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
}

// LoginRequest for authentication
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest for new users
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// BookingRequest for journey bookings
type BookingRequest struct {
	JourneyID string `json:"journeyId" binding:"required"`
	UserID    string `json:"userId" binding:"required"`
}
