package model

import "time"

type PartnerContractRule struct {
	ID                    int       `json:"id"`
	RevenueSharingPercent float64   `json:"revenue_sharing_percent"`
	MaxWarningCount       int       `json:"max_warning_count"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
