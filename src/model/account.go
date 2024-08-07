package model

import "time"

type AccountStatus string

const (
	AccountStatusWaitingConfirmPhoneNumber AccountStatus = "waiting_confirm_phone_number"
	AccountStatusActive                    AccountStatus = "active"
	AccountStatusInactive                  AccountStatus = "inactive"
	AccountStatusNoFilter                  AccountStatus = "no_filter"
)

type Account struct {
	ID                       int           `json:"id"`
	RoleID                   RoleID        `json:"role_id"`
	Role                     *Role         `json:"role,omitempty"`
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
	BankNumber               string        `json:"bank_number"`
	BankOwner                string        `json:"bank_owner"`
	BankName                 string        `json:"bank_name"`
	QRCodeURL                string        `json:"qr_code_url"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}
