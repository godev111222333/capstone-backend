package model

import "time"

type (
	DocumentCategory string
	DocumentStatus   string
)

const (
	ExtensionPDF = "pdf"
)

const (
	DocumentCategoryCarImages          DocumentCategory = "CAR_IMAGES"
	DocumentCategoryCaveat             DocumentCategory = "CAR_CAVEAT"
	DocumentCategoryDrivingLicense     DocumentCategory = "DRIVING_LICENSE"
	DocumentCategoryQRCodeImage        DocumentCategory = "QR_CODE"
	DocumentCategoryAvatarImage        DocumentCategory = "AVATAR"
	DocumentCategoryPrepayQRCodeImage  DocumentCategory = "PREPAY_QR_CODE"
	DocumentCategoryCollateralAssets   DocumentCategory = "COLLATERAL_ASSETS"
	DocumentCategoryReceivingCarImages DocumentCategory = "RECEIVING_CAR_IMAGES"
	DocumentStatusActive               DocumentStatus   = "active"
	DocumentStatusInactive             DocumentStatus   = "inactive"
)

type Document struct {
	ID        int              `json:"id"`
	AccountID int              `json:"account_id"`
	Account   Account          `json:"account"`
	Url       string           `json:"url"`
	Extension string           `json:"extension"`
	Category  DocumentCategory `json:"category"`
	Status    DocumentStatus   `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
