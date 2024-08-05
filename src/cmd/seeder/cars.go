package seeder

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/api"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type Car struct {
	ID                    int                         `csv:"id"`
	PartnerID             int                         `csv:"partner_id"`
	CarModelID            int                         `csv:"car_model_id"`
	LicensePlate          string                      `csv:"license_plate"`
	ParkingLot            model.ParkingLot            `csv:"parking_lot"`
	Description           string                      `csv:"description"`
	Fuel                  model.Fuel                  `csv:"fuel"`
	Motion                model.Motion                `csv:"motion"`
	Price                 int                         `csv:"price"`
	Status                model.CarStatus             `csv:"status"`
	PartnerContractRuleID int                         `csv:"partner_contract_rule_id"`
	BankName              string                      `csv:"bank_name"`
	BankNumber            string                      `csv:"bank_number"`
	BankOwner             string                      `csv:"bank_owner"`
	StartDate             DateTime                    `csv:"start_date"`
	EndDate               DateTime                    `csv:"end_date"`
	Period                int                         `csv:"period"`
	PartnerContractUrl    string                      `csv:"partner_contract_url"`
	PartnerContractStatus model.PartnerContractStatus `csv:"partner_contract_status"`
	WarningCount          int                         `csv:"warning_count"`
	CreatedAt             DateTime                    `csv:"created_at"`
	UpdatedAt             DateTime                    `csv:"updated_at"`
}

func (c *Car) ToDbCar() *model.Car {
	return &model.Car{
		ID:                    c.ID,
		PartnerID:             c.PartnerID,
		CarModelID:            c.CarModelID,
		LicensePlate:          c.LicensePlate,
		ParkingLot:            c.ParkingLot,
		Description:           c.Description,
		Fuel:                  c.Fuel,
		Motion:                c.Motion,
		Price:                 c.Price,
		Status:                c.Status,
		PartnerContractRuleID: c.PartnerContractRuleID,
		BankName:              c.BankName,
		BankNumber:            c.BankNumber,
		BankOwner:             c.BankOwner,
		StartDate:             c.StartDate.Time,
		EndDate:               c.EndDate.Time,
		Period:                c.Period,
		PartnerContractUrl:    c.PartnerContractUrl,
		PartnerContractStatus: c.PartnerContractStatus,
		WarningCount:          c.WarningCount,
		CreatedAt:             c.CreatedAt.Time,
		UpdatedAt:             c.CreatedAt.Time,
	}
}

func SeedCars(server *api.Server, dbStore *store.DbStore) error {
	cars := make([]*Car, 0)
	accountFile, err := os.OpenFile(toFilePath(CarsFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer accountFile.Close()

	if err := gocsv.UnmarshalFile(accountFile, &cars); err != nil {
		return err
	}

	dbCars := make([]*model.Car, len(cars))
	for i, a := range cars {
		dbCars[i] = a.ToDbCar()
	}

	if err := dbStore.CarStore.CreateBatch(dbCars); err != nil {
		return err
	}

	for _, car := range dbCars {
		partner, err := dbStore.AccountStore.GetByID(car.PartnerID)
		if err != nil {
			return err
		}

		now := car.CreatedAt.Add(time.Hour)
		if err := server.InternalRenderPartnerContractPDF(partner, car, now); err != nil {
			return err
		}
	}

	return nil
}
