package model

import "time"

type AccountStatus string

const (
	AccountStatusEnable  AccountStatus = "enable"
	AccountStatusDisable AccountStatus = "disable"
)

type Account struct {
	ID                       int           `json:"id,omitempty"`
	RoleID                   RoleID        `json:"role_id,omitempty"`
	FirstName                string        `json:"first_name,omitempty"`
	LastName                 string        `json:"last_name,omitempty"`
	PhoneNumber              string        `json:"phone_number,omitempty"`
	Email                    string        `json:"email,omitempty"`
	IdentificationCardNumber string        `json:"identification_card_number,omitempty"`
	DateOfBirth              time.Time     `json:"date_of_birth"`
	Password                 string        `json:"password,omitempty"`
	AvatarURL                string        `json:"avatar_url,omitempty"`
	Status                   AccountStatus `json:"status,omitempty"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}
