package model

import "time"

type OTPType string

type OTPStatus string

const (
	OTPTypeRegister OTPType = "Register"
)

const (
	OTPStatusVerified OTPStatus = "Verified"
	OTPStatusSent     OTPStatus = "Sent"
)

type OTP struct {
	ID           int
	AccountEmail string
	OTP          string
	Status       OTPStatus
	OtpType      OTPType
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
