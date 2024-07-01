package model

import "time"

type (
	CustomerContractImageCategory string
	CustomerContractImageStatus   string
)

const (
	CustomerContractImageCategoryCollateralAssets   CustomerContractImageCategory = "COLLATERAL_ASSETS"
	CustomerContractImageCategoryReceivingCarImages CustomerContractImageCategory = "RECEIVING_CAR_IMAGES"

	CustomerContractImageStatusActive   CustomerContractImageStatus = "active"
	CustomerContractImageStatusInactive CustomerContractImageStatus = "inactive"
)

type CustomerContractImage struct {
	ID                 int                           `json:"id"`
	CustomerContractID int                           `json:"customer_contract_id"`
	CustomerContract   *CustomerContract             `json:"customer_contract,omitempty"`
	URL                string                        `json:"url"`
	Category           CustomerContractImageCategory `json:"category"`
	Status             CustomerContractImageStatus   `json:"status"`
	CreatedAt          time.Time                     `json:"created_at"`
	UpdatedAt          time.Time                     `json:"updated_at"`
}
