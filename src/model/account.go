package model

import "time"

type AccountStatus string

const (
	AccountStatusWaitingConfirmEmail AccountStatus = "waiting_confirm_email"
	AccountStatusActive              AccountStatus = "active"
	AccountStatusInactive            AccountStatus = "inactive"
)

type Account struct {
	ID                       int           `json:"id"`
	RoleID                   RoleID        `json:"role_id"`
	Role                     Role          `json:"role"`
	FirstName                string        `json:"first_name"`
	LastName                 string        `json:"last_name"`
	PhoneNumber              string        `json:"phone_number"`
	Email                    string        `json:"email"`
	IdentificationCardNumber string        `json:"identification_card_number"`
	Password                 string        `json:"password"`
	AvatarURL                string        `json:"avatar_url"`
	DrivingLicense           string        `json:"driving_license"`
	Status                   AccountStatus `json:"status"`
	DateOfBirth              time.Time     `json:"date_of_birth"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}
