package model

import "time"

type CarDocument struct {
	ID         int       `json:"id"`
	CarID      int       `json:"car_id"`
	Car        Car       `json:"car"`
	DocumentID int       `json:"document_id"`
	Document   Document  `json:"document"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
