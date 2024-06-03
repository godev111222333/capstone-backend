package model

import "time"

type CarModel struct {
	ID            int       `json:"id"`
	Brand         string    `json:"brand"`
	Model         string    `json:"model"`
	Year          int       `json:"year"`
	NumberOfSeats int       `json:"number_of_seats"`
	BasedPrice    int       `json:"based_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
