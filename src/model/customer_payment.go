package model

import "time"

type (
	PaymentType   string
	PaymentStatus string
)

const (
	PaymentTypePrePay       PaymentType   = "pre_pay"
	PaymentTypeRemainingPay PaymentType   = "remaining_pay"
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusPaid       PaymentStatus = "paid"
)

type CustomerPayment struct {
	ID                 int              `json:"id"`
	CustomerContractID int              `json:"customer_contract_id"`
	CustomerContract   CustomerContract `json:"customer_contract" gorm:"foreignKey:CustomerContractID"`
	PaymentType        PaymentType      `json:"payment_type"`
	Amount             int              `json:"amount"`
	Note               string           `json:"note"`
	Status             PaymentStatus    `json:"status"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}
