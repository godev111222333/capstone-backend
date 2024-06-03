package model

import "time"

type (
	ParkingLot string
	Fuel       string
	CarStatus  string
	Motion     string
)

const (
	ParkingLotHome ParkingLot = "home"

	ParkingLotGarage ParkingLot = "garage"

	FuelGas                     Fuel      = "gas"
	FuelOil                     Fuel      = "oil"
	FuelElectricity             Fuel      = "electricity"
	MotionAutomaticTransmission Motion    = "automatic_transmission"
	MotionManualTransmission    Motion    = "manual_transmission"
	CarStatusPendingApproval    CarStatus = "pending_approval"
	CarStatusApproved           CarStatus = "approved"
	CarStatusRejected           CarStatus = "rejected"
	CarStatusActive             CarStatus = "active"
	CarStatusWaitingDelivery    CarStatus = "waiting_car_delivery"
)

type Car struct {
	ID           int        `json:"id"`
	AccountID    int        `json:"account_id"`
	Account      Account    `json:"account"`
	CarModelID   int        `json:"car_model_id"`
	CarModel     CarModel   `json:"car_model"`
	LicensePlate string     `json:"license_plate"`
	ParkingLot   ParkingLot `json:"parking_lot"`
	Description  string     `json:"description"`
	Fuel         Fuel       `json:"fuel"`
	Motion       Motion     `json:"motion"`
	Price        int        `json:"price"`
	Status       CarStatus  `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
