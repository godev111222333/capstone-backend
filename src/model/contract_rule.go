package model

import "time"

type ContractRule struct {
	ID                    int       `json:"id"`
	InsurancePercent      float64   `json:"insurance_percent"`
	PrepayPercent         float64   `json:"prepay_percent"`
	RevenueSharingPercent float64   `json:"revenue_sharing_percent"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
