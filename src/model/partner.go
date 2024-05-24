package model

import "time"

type Partner struct {
	ID        int       `json:"id,omitempty"`
	AccountID int       `json:"account_id,omitempty"`
	Account   Account   `json:"account"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
