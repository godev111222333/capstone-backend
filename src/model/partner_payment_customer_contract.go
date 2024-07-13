package model

import "time"

type PartnerPaymentCustomerContract struct {
	ID                      int                    `json:"id"`
	PartnerPaymentHistoryID int                    `json:"partner_payment_history_id,omitempty"`
	PartnerPaymentHistory   *PartnerPaymentHistory `json:"partner_payment_history"`
	CustomerContractID      int                    `json:"customer_contract_id"`
	CustomerContract        *CustomerContract      `json:"customer_contract,omitempty"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
}
