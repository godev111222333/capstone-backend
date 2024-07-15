package model

import "time"

type PartnerPaymentHistoryStatus string

const (
	PartnerPaymentHistoryStatusPending  PartnerPaymentHistoryStatus = "pending"
	PartnerPaymentHistoryStatusPaid     PartnerPaymentHistoryStatus = "paid"
	PartnerPaymentHistoryStatusNoFilter PartnerPaymentHistoryStatus = "no_filter"
)

type PartnerPaymentHistory struct {
	ID         int                         `json:"id"`
	PartnerID  int                         `json:"partner_id"`
	Partner    *Account                    `json:"partner,omitempty" gorm:"foreignKey:PartnerID"`
	StartDate  time.Time                   `json:"start_date"`
	EndDate    time.Time                   `json:"end_date"`
	Amount     int                         `json:"amount"`
	PaymentURL string                      `json:"payment_url"`
	Status     PartnerPaymentHistoryStatus `json:"status"`
	CreatedAt  time.Time                   `json:"created_at"`
	UpdatedAt  time.Time                   `json:"updated_at"`
}
