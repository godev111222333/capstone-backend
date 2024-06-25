package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	PhoneNumber  string    `json:"phone_number"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
