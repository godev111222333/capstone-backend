package model

import "time"

type (
	DocumentCategory string
	DocumentStatus   string
)

const (
	DocumentCategoryCarImages   DocumentCategory = "CAR_IMAGES"
	DocumentCategoryCaveat      DocumentCategory = "CAR_CAVEAT"
	DocumentCategoryQRCodeImage DocumentCategory = "QR_CODE"
	DocumentCategoryAvatarImage DocumentCategory = "AVATAR"
	DocumentStatusActive        DocumentStatus   = "active"
	DocumentStatusInactive      DocumentStatus   = "inactive"
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
