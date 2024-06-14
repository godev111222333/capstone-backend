package model

import "time"

type PartnerContractStatus string

const (
	PartnerContractStatusWaitingForAgreement PartnerContractStatus = "waiting_for_agreement"
	PartnerContractStatusAgreed              PartnerContractStatus = "agreed"
)

type PartnerContract struct {
	ID        int                   `json:"id"`
	CarID     int                   `json:"car_id"`
	Car       Car                   `json:"car"`
	StartDate time.Time             `json:"start_date"`
	EndDate   time.Time             `json:"end_date"`
	Url       string                `json:"url"`
	Status    PartnerContractStatus `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
