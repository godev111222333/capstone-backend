package model

import "time"

type GarageConfigType string

const (
	GarageConfigTypeMax4Seats  GarageConfigType = "MAX_4_SEATS"
	GarageConfigTypeMax7Seats  GarageConfigType = "MAX_7_SEATS"
	GarageConfigTypeMax15Seats GarageConfigType = "MAX_15_SEATS"
)

type GarageConfig struct {
	ID        int              `json:"id"`
	Type      GarageConfigType `json:"type"`
	Maximum   int              `json:"maximum"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
