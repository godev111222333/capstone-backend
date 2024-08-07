package model

import "time"

type (
	CustomerContractStatus string
	CollateralType         string
	FeedBackStatus         string
)

const (
	CustomerContractStatusWaitingContractAgreement CustomerContractStatus = "waiting_for_agreement"
	CustomerContractStatusWaitingPartnerApproval   CustomerContractStatus = "waiting_partner_approval"
	CustomerContractStatusWaitingContractPayment   CustomerContractStatus = "waiting_contract_payment"
	CustomerContractStatusOrdered                  CustomerContractStatus = "ordered"
	CustomerContractStatusRenting                  CustomerContractStatus = "renting"
	CustomerContractStatusCompleted                CustomerContractStatus = "completed"
	CustomerContractStatusCancel                   CustomerContractStatus = "canceled"
	CustomerContractStatusNoFilter                 CustomerContractStatus = "no_filter"

	FeedbackStatusActive   FeedBackStatus = "active"
	FeedbackStatusInactive FeedBackStatus = "inactive"

	CollateralTypeCash      CollateralType = "cash"
	CollateralTypeMotorbike CollateralType = "motorbike"
)

type CustomerContract struct {
	ID                      int                    `json:"id"`
	CustomerID              int                    `json:"customer_id"`
	Customer                *Account               `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	CarID                   int                    `json:"car_id"`
	Car                     Car                    `json:"car"`
	StartDate               time.Time              `json:"start_date"`
	EndDate                 time.Time              `json:"end_date"`
	Status                  CustomerContractStatus `json:"status"`
	Reason                  string                 `json:"reason"`
	RentPrice               int                    `json:"rent_price"`
	InsuranceAmount         int                    `json:"insurance_amount"`
	CollateralType          CollateralType         `json:"collateral_type"`
	IsReturnCollateralAsset bool                   `json:"is_return_collateral_asset"`
	Url                     string                 `json:"url"`
	BankName                string                 `json:"bank_name"`
	BankNumber              string                 `json:"bank_number"`
	BankOwner               string                 `json:"bank_owner"`
	CustomerContractRuleID  int                    `json:"customer_contract_rule_id"`
	CustomerContractRule    CustomerContractRule   `gorm:"foreignKey:CustomerContractRuleID" json:"customer_contract_rule,omitempty"`
	FeedbackContent         string                 `json:"feedback_content"`
	FeedbackRating          int                    `json:"feedback_rating"`
	FeedbackStatus          FeedBackStatus         `json:"feedback_status"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
}
