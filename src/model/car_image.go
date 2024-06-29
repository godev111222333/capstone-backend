package model

import "time"

type (
	CarImageCategory string
	CarImageStatus   string
)

const (
	ExtensionPDF = ".pdf"
)

const (
	CarImageCategoryImages CarImageCategory = "CAR_IMAGES"
	CarImageCategoryCaveat CarImageCategory = "CAR_CAVEAT"
	CarImageStatusActive   CarImageStatus   = "active"
	CarImageStatusInactive CarImageStatus   = "inactive"
)

type CarImage struct {
	ID        int              `json:"id"`
	CarID     int              `json:"car_id"`
	Car       Car              `json:"car"`
	URL       string           `json:"url"`
	Category  CarImageCategory `json:"category"`
	Status    CarImageStatus   `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
