package model

import "time"

type (
	PaymentType   string
	PaymentStatus string
)

const (
	PaymentTypePrePay               PaymentType   = "pre_pay"
	PaymentTypeRemainingPay         PaymentType   = "remaining_pay"
	PaymentTypeCollateralCash       PaymentType   = "collateral_cash"
	PaymentTypeReturnCollateralCash PaymentType   = "return_collateral_cash"
	PaymentTypeOther                PaymentType   = "other"
	PaymentStatusPending            PaymentStatus = "pending"
	PaymentStatusPaid               PaymentStatus = "paid"
	PaymentStatusCanceled           PaymentStatus = "canceled"
	PaymentStatusNoFilter           PaymentStatus = "no_filter"
)

type CustomerPayment struct {
	ID                 int               `json:"id"`
	CustomerContractID int               `json:"customer_contract_id"`
	CustomerContract   *CustomerContract `json:"customer_contract,omitempty" gorm:"foreignKey:CustomerContractID"`
	PaymentType        PaymentType       `json:"payment_type"`
	Amount             int               `json:"amount"`
	Note               string            `json:"note"`
	Status             PaymentStatus     `json:"status"`
	PaymentURL         string            `json:"payment_url"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}
