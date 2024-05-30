package model

import "time"

type AccountStatus string

const (
	AccountStatusWaitingConfirmEmail AccountStatus = "waiting_confirm_email"
	AccountStatusEnable              AccountStatus = "active"
	AccountStatusDisable             AccountStatus = "inactive"
)

type Account struct {
	ID                       int           `json:"id,omitempty"`
	RoleID                   RoleID        `json:"role_id,omitempty"`
	Role                     Role          `json:"role"`
	FirstName                string        `json:"first_name,omitempty"`
	LastName                 string        `json:"last_name,omitempty"`
	PhoneNumber              string        `json:"phone_number,omitempty"`
	Email                    string        `json:"email,omitempty"`
	IdentificationCardNumber string        `json:"identification_card_number,omitempty"`
	Password                 string        `json:"password,omitempty"`
	AvatarURL                string        `json:"avatar_url,omitempty"`
	DrivingLicense           string        `json:"driving_license"`
	Status                   AccountStatus `json:"status,omitempty"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}
