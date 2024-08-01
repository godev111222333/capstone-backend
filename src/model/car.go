package model

import "time"

type (
	ParkingLot string
	Fuel       string
	CarStatus  string
	Motion     string
)

type PartnerContractStatus string

const (
	PartnerContractStatusWaitingForApproval  PartnerContractStatus = "waiting_for_approval"
	PartnerContractStatusWaitingForAgreement PartnerContractStatus = "waiting_for_agreement"
	PartnerContractStatusAgreed              PartnerContractStatus = "agreed"
)

const (
	ParkingLotHome ParkingLot = "home"

	ParkingLotGarage ParkingLot = "garage"

	FuelGas                                     Fuel      = "gas"
	FuelOil                                     Fuel      = "oil"
	FuelElectricity                             Fuel      = "electricity"
	MotionAutomaticTransmission                 Motion    = "automatic_transmission"
	MotionManualTransmission                    Motion    = "manual_transmission"
	CarStatusPendingApplication                 CarStatus = "pending_application"
	CarStatusPendingApplicationPendingCarImages           = CarStatusPendingApplication + ":pending_car_images"
	CarStatusPendingApplicationPendingCarCaveat           = CarStatusPendingApplication + ":pending_car_caveat"
	CarStatusPendingApplicationPendingPrice               = CarStatusPendingApplication + ":pending_price"
	CarStatusPendingApproval                    CarStatus = "pending_approval"
	CarStatusApproved                           CarStatus = "approved"
	CarStatusRejected                           CarStatus = "rejected"
	CarStatusActive                             CarStatus = "active"
	CarStatusInactive                           CarStatus = "inactive"
	CarStatusWaitingDelivery                    CarStatus = "waiting_car_delivery"
	CarStatusNoFilter                           CarStatus = "no_filter"
)

type Car struct {
	ID                    int                   `json:"id"`
	PartnerID             int                   `json:"partner_id"`
	Account               *Account              `json:"account,omitempty" gorm:"foreignKey:PartnerID"`
	CarModelID            int                   `json:"car_model_id"`
	CarModel              CarModel              `json:"car_model,omitempty"`
	LicensePlate          string                `json:"license_plate"`
	ParkingLot            ParkingLot            `json:"parking_lot"`
	Description           string                `json:"description"`
	Fuel                  Fuel                  `json:"fuel"`
	Motion                Motion                `json:"motion"`
	Price                 int                   `json:"price"`
	Status                CarStatus             `json:"status"`
	PartnerContractRuleID int                   `json:"partner_contract_rule_id"`
	PartnerContractRule   PartnerContractRule   `json:"partner_contract_rule" gorm:"foreignKey:PartnerContractRuleID"`
	BankName              string                `json:"bank_name"`
	BankNumber            string                `json:"bank_number"`
	BankOwner             string                `json:"bank_owner"`
	StartDate             time.Time             `json:"start_date"`
	EndDate               time.Time             `json:"end_date"`
	Period                int                   `json:"period"`
	PartnerContractUrl    string                `json:"partner_contract_url"`
	PartnerContractStatus PartnerContractStatus `json:"partner_contract_status"`
	WarningCount          int                   `json:"warning_count"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
}

type PartnerContract struct {
	CarID                 int                   `json:"car_id"`
	RevenueSharingPercent float64               `json:"revenue_sharing_percent"`
	MaxWarningCount       int                   `json:"max_warning_count"`
	BankName              string                `json:"bank_name"`
	BankNumber            string                `json:"bank_number"`
	BankOwner             string                `json:"bank_owner"`
	StartDate             time.Time             `json:"start_date"`
	EndDate               time.Time             `json:"end_date"`
	Url                   string                `json:"url"`
	Status                PartnerContractStatus `json:"status"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
}

func (c *Car) ToPartnerContract() *PartnerContract {
	return &PartnerContract{
		CarID:                 c.ID,
		RevenueSharingPercent: c.PartnerContractRule.RevenueSharingPercent,
		MaxWarningCount:       c.PartnerContractRule.MaxWarningCount,
		BankName:              c.BankName,
		BankNumber:            c.BankNumber,
		BankOwner:             c.BankOwner,
		StartDate:             c.StartDate,
		EndDate:               c.EndDate,
		Url:                   c.PartnerContractUrl,
		Status:                c.PartnerContractStatus,
		CreatedAt:             c.CreatedAt,
		UpdatedAt:             c.UpdatedAt,
	}
}

type CarJoinCarModel struct {
	CarID                 int                 `json:"car_id"`
	CarModelID            int                 `json:"car_model_id"`
	CarModel              CarModel            `json:"car_model" gorm:"foreignKey:CarModelID"`
	PartnerID             int                 `json:"partner_id"`
	Brand                 string              `json:"brand"`
	Model                 string              `json:"model"`
	Year                  int                 `json:"year"`
	NumberOfSeats         int                 `json:"number_of_seats"`
	LicensePlate          string              `json:"license_plate"`
	ParkingLot            ParkingLot          `json:"parking_lot"`
	Description           string              `json:"description"`
	Fuel                  Fuel                `json:"fuel"`
	Motion                Motion              `json:"motion"`
	Price                 int                 `json:"price"`
	Status                CarStatus           `json:"status"`
	Period                int                 `json:"period"`
	PartnerContractRuleID int                 `json:"partner_contract_rule_id"`
	PartnerContractRule   PartnerContractRule `json:"partner_contract_rule" gorm:"foreignKey:PartnerContractRuleID"`
	CreatedAt             time.Time           `json:"created_at"`
	UpdatedAt             time.Time           `json:"updated_at"`
}

func (m *CarJoinCarModel) ToCar() *Car {
	return &Car{
		ID:                    m.CarID,
		PartnerID:             m.PartnerID,
		CarModelID:            m.CarModelID,
		CarModel:              m.CarModel,
		LicensePlate:          m.LicensePlate,
		ParkingLot:            m.ParkingLot,
		Description:           m.Description,
		Fuel:                  m.Fuel,
		Motion:                m.Motion,
		Price:                 m.Price,
		Status:                m.Status,
		PartnerContractRule:   m.PartnerContractRule,
		PartnerContractRuleID: m.PartnerContractRuleID,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}

var CarStates = []CarStatus{
	CarStatusPendingApplicationPendingCarImages,
	CarStatusPendingApplicationPendingCarCaveat,
	CarStatusPendingApplicationPendingPrice,
	CarStatusPendingApproval,
	CarStatusWaitingDelivery,
	CarStatusActive,
}

func MoveNextCarState(curState CarStatus) CarStatus {
	for i := 0; i < len(CarStates); i++ {
		if curState == CarStates[i] {
			return CarStates[i+1]
		}
	}

	return CarStatusRejected
}
