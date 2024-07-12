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
	PhoneNumber string    `json:"phone_number"`
	OTP         string    `json:"OTP"`
	Status      OTPStatus `json:"status"`
	OtpType     OTPType   `json:"otp_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}
