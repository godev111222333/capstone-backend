package model

import "time"

type CustomerContractDocument struct {
	ID                 int              `json:"id"`
	CustomerContractID int              `json:"customer_contract_id"`
	CustomerContract   CustomerContract `json:"customer_contract"`
	DocumentID         int              `json:"document_id"`
	Document           Document         `json:"document"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}
