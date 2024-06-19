package model

import "time"

type CustomerPaymentDocument struct {
	ID                int             `json:"id"`
	CustomerPaymentID int             `json:"customer_payment_id"`
	CustomerPayment   CustomerPayment `json:"customer_payment" gorm:"foreignKey:CustomerPaymentID"`
	DocumentID        int             `json:"document_id"`
	Document          Document        `json:"document" gorm:"foreignKey:DocumentID"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
