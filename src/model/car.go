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
	CarStatusNoFilter           CarStatus = "no_filter"
)

type Car struct {
	ID           int        `json:"id"`
	PartnerID    int        `json:"partner_id"`
	Account      Account    `json:"account,omitempty" gorm:"foreignKey:PartnerID"`
	CarModelID   int        `json:"car_model_id"`
	CarModel     CarModel   `json:"car_model,omitempty"`
	LicensePlate string     `json:"license_plate"`
	ParkingLot   ParkingLot `json:"parking_lot"`
	Description  string     `json:"description"`
	Fuel         Fuel       `json:"fuel"`
	Motion       Motion     `json:"motion"`
	Price        int        `json:"price"`
	Status       CarStatus  `json:"status"`
	Period       int        `json:"period"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CarJoinCarModel struct {
	ID            int        `json:"id"`
	CarModelID    int        `json:"car_model_id"`
	PartnerID     int        `json:"partner_id"`
	Brand         string     `json:"brand"`
	Model         string     `json:"model"`
	Year          int        `json:"year"`
	NumberOfSeats int        `json:"number_of_seats"`
	LicensePlate  string     `json:"license_plate"`
	ParkingLot    ParkingLot `json:"parking_lot"`
	Description   string     `json:"description"`
	Fuel          Fuel       `json:"fuel"`
	Motion        Motion     `json:"motion"`
	Price         int        `json:"price"`
	Status        CarStatus  `json:"status"`
	Period        int        `json:"period"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (m *CarJoinCarModel) ToCar() *Car {
	return &Car{
		ID:           m.ID,
		PartnerID:    m.PartnerID,
		CarModelID:   m.CarModelID,
		LicensePlate: m.LicensePlate,
		ParkingLot:   m.ParkingLot,
		Description:  m.Description,
		Fuel:         m.Fuel,
		Motion:       m.Motion,
		Price:        m.Price,
		Status:       m.Status,
		Period:       m.Period,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
