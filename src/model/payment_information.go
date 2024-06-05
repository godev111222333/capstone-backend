package model

import "time"

type PaymentInformation struct {
	ID         int       `json:"id"`
	AccountID  int       `json:"account_id"`
	Account    Account   `json:"account"`
	BankNumber string    `json:"bank_number"`
	BankOwner  string    `json:"bank_owner"`
	BankName   string    `json:"bank_name"`
	QRCodeURL  string    `json:"qr_code_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
