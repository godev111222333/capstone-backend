package model

import "time"

type DrivingLicenseImageStatus string

const (
	DrivingLicenseImageStatusActive   DrivingLicenseImageStatus = "active"
	DrivingLicenseImageStatusInactive DrivingLicenseImageStatus = "inactive"
)

type DrivingLicenseImage struct {
	ID        int                       `json:"id"`
	AccountID int                       `json:"account_id"`
	Account   Account                   `json:"account"`
	URL       string                    `json:"url"`
	Status    DrivingLicenseImageStatus `json:"status"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}
