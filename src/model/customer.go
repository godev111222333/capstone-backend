package model

import "time"

type Customer struct {
	ID             int       `json:"id,omitempty"`
	AccountID      int       `json:"account_id,omitempty"`
	Account        Account   `json:"account"`
	DrivingLicense string    `json:"driving_license,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
