package model

import "time"

type CarModel struct {
	ID            int
	Brand         string
	Model         string
	Year          int
	NumberOfSeats int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
