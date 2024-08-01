package model

import "time"

type CustomerContractRule struct {
	ID                   int       `json:"id"`
	InsurancePercent     float64   `json:"insurance_percent"`
	PrepayPercent        float64   `json:"prepay_percent"`
	CollateralCashAmount int       `json:"collateral_cash_amount"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
