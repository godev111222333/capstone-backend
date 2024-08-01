package store

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestCarStore(t *testing.T) {
	t.Run("create car successfully", func(t *testing.T) {
		carModel := &model.CarModel{
			Brand:         "BMW",
			Model:         "X8",
			Year:          2024,
			NumberOfSeats: 4,
			BasedPrice:    300_000,
		}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
		partner := &model.Account{
			RoleID:    model.RoleIDPartner,
			FirstName: "Son Le",
			Status:    model.AccountStatusActive,
		}
		require.NoError(t, TestDb.AccountStore.Create(partner))
		car := &model.Car{
			PartnerID:             partner.ID,
			CarModelID:            carModel.ID,
			LicensePlate:          "7777",
			ParkingLot:            model.ParkingLotHome,
			Description:           "Beautiful car",
			Fuel:                  model.FuelElectricity,
			Motion:                model.MotionAutomaticTransmission,
			Price:                 550_000,
			Status:                model.CarStatusActive,
			PartnerContractRuleID: 1,
		}
		require.NoError(t, TestDb.CarStore.Create(car))
	})

	t.Run("get owned car successfully", func(t *testing.T) {
		partner := &model.Account{
			RoleID:      model.RoleIDPartner,
			PhoneNumber: "1111",
			FirstName:   "Cuong dola",
			Status:      model.AccountStatusActive,
		}
		require.NoError(t, TestDb.AccountStore.Create(partner))
		carModel := &model.CarModel{
			Brand: "Bugatti",
		}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))

		for i := 1; i <= 2; i++ {
			car := &model.Car{
				PartnerID:             partner.ID,
				CarModelID:            carModel.ID,
				LicensePlate:          "86A" + strconv.Itoa(i),
				Status:                model.CarStatusActive,
				PartnerContractRuleID: 1,
			}
			require.NoError(t, TestDb.CarStore.Create(car))
		}
		cars := []*model.Car{
			{
				PartnerID:             partner.ID,
				CarModelID:            carModel.ID,
				LicensePlate:          "86AX1",
				Status:                model.CarStatusPendingApproval,
				PartnerContractRuleID: 1,
			},
			{
				PartnerID:             partner.ID,
				CarModelID:            carModel.ID,
				LicensePlate:          "86AX2",
				Status:                model.CarStatusPendingApplicationPendingCarImages,
				PartnerContractRuleID: 1,
			},
		}
		for _, car := range cars {
			require.NoError(t, TestDb.CarStore.Create(car))
		}

		cars, err := TestDb.CarStore.GetByPartner(partner.ID, 0, 0, model.CarStatusNoFilter)
		require.NoError(t, err)
		require.Len(t, cars, 4)
		cars, err = TestDb.CarStore.GetByPartner(partner.ID, 0, 2, model.CarStatusPendingApproval)
		require.NoError(t, err)
		require.Len(t, cars, 1)
		cars, err = TestDb.CarStore.GetByPartner(partner.ID, 0, 2, model.CarStatusPendingApplication)
		require.NoError(t, err)
		require.Len(t, cars, 1)
	})

	t.Run("get all", func(t *testing.T) {
		partner := &model.Account{
			RoleID:      model.RoleIDPartner,
			PhoneNumber: "222222",
			FirstName:   "Cuong dola 2",
			Status:      model.AccountStatusActive,
		}
		require.NoError(t, TestDb.AccountStore.Create(partner))
		carModel := &model.CarModel{
			Brand: "Bugatti",
		}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))

		for i := 1; i <= 5; i++ {
			car := &model.Car{
				PartnerID:             partner.ID,
				CarModelID:            carModel.ID,
				LicensePlate:          "51A" + strconv.Itoa(i),
				Status:                model.CarStatusActive,
				PartnerContractRuleID: 1,
			}
			require.NoError(t, TestDb.CarStore.Create(car))
		}

		for i := 1; i <= 3; i++ {
			car := &model.Car{
				PartnerID:             partner.ID,
				CarModelID:            carModel.ID,
				LicensePlate:          "xxxx-" + strconv.Itoa(i),
				Status:                model.CarStatusPendingApproval,
				PartnerContractRuleID: 1,
			}
			require.NoError(t, TestDb.CarStore.Create(car))
		}

		cars, err := TestDb.CarStore.SearchCars(0, 100, model.CarStatusNoFilter, "")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(cars), 8)

		cars, err = TestDb.CarStore.SearchCars(0, 100, model.CarStatusPendingApproval, "")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(cars), 3)

		cars, err = TestDb.CarStore.SearchCars(0, 1, model.CarStatusPendingApproval, "")
		require.NoError(t, err)
		require.Len(t, cars, 1)

		cars, err = TestDb.CarStore.SearchCars(0, 3, model.CarStatusPendingApproval, "Cuong dola 2")
		require.NoError(t, err)
		require.Len(t, cars, 3)
		require.NotEmpty(t, cars[0].CarModel.Brand)
	})

	t.Run("find cars", func(t *testing.T) {
		require.NoError(t, ResetDb(TestConfig.Database))
		carModels := []*model.CarModel{
			{Brand: "toyota", Model: "X8", Year: 2024, NumberOfSeats: 4},
			{Brand: "mec", Model: "maybach s450", Year: 2024, NumberOfSeats: 7},
			{Brand: "audi", Model: "G9", Year: 2024, NumberOfSeats: 15},
		}
		require.NoError(t, TestDb.CarModelStore.Create(carModels))
		partner1 := &model.Account{PhoneNumber: "3333", Password: "p1", RoleID: model.RoleIDPartner}
		partner2 := &model.Account{PhoneNumber: "4444", Password: "p2", RoleID: model.RoleIDPartner}
		require.NoError(t, TestDb.AccountStore.Create(partner1))
		require.NoError(t, TestDb.AccountStore.Create(partner2))

		now := time.Now()

		cars := []*model.Car{
			{
				PartnerID:             partner1.ID,
				CarModelID:            carModels[0].ID,
				LicensePlate:          "86B1",
				ParkingLot:            model.ParkingLotGarage,
				Fuel:                  model.FuelElectricity,
				Motion:                model.MotionAutomaticTransmission,
				Price:                 100_000,
				Status:                model.CarStatusActive,
				EndDate:               now.AddDate(1, 0, 0),
				PartnerContractRuleID: 1,
			},
			{
				PartnerID:             partner2.ID,
				CarModelID:            carModels[1].ID,
				LicensePlate:          "86B2",
				ParkingLot:            model.ParkingLotHome,
				Fuel:                  model.FuelElectricity,
				Motion:                model.MotionAutomaticTransmission,
				Price:                 200_000,
				Status:                model.CarStatusActive,
				EndDate:               now.AddDate(1, 0, 0),
				PartnerContractRuleID: 1,
			},
			{
				PartnerID:             partner2.ID,
				CarModelID:            carModels[2].ID,
				LicensePlate:          "86B3",
				ParkingLot:            model.ParkingLotGarage,
				Fuel:                  model.FuelGas,
				Motion:                model.MotionManualTransmission,
				Price:                 300_000,
				Status:                model.CarStatusActive,
				EndDate:               now.AddDate(1, 0, 0),
				PartnerContractRuleID: 1,
			},
		}
		for _, car := range cars {
			require.NoError(t, TestDb.CarStore.Create(car))
		}

		customer := &model.Account{Email: "c1@gmail.com", Password: "c1", RoleID: model.RoleIDCustomer}
		require.NoError(t, TestDb.AccountStore.Create(customer))

		contractRuleID := 1
		customerContracts := []*model.CustomerContract{
			{
				CustomerID:             customer.ID,
				CarID:                  cars[0].CarModelID,
				StartDate:              now.Add(10 * time.Hour),
				EndDate:                now.Add(10 * time.Hour * time.Duration(3)),
				Status:                 model.CustomerContractStatusOrdered,
				CustomerContractRuleID: contractRuleID,
			},
			{
				CustomerID:             customer.ID,
				CarID:                  cars[0].CarModelID,
				StartDate:              now.Add(10 * time.Hour * time.Duration(7)),
				EndDate:                now.Add(10 * time.Hour * time.Duration(9)),
				Status:                 model.CustomerContractStatusOrdered,
				CustomerContractRuleID: contractRuleID,
			},
			{
				CustomerID:             customer.ID,
				CarID:                  cars[0].CarModelID,
				StartDate:              now.Add(10 * time.Hour * time.Duration(15)),
				EndDate:                now.Add(10 * time.Hour * time.Duration(20)),
				Status:                 model.CustomerContractStatusWaitingContractPayment,
				CustomerContractRuleID: contractRuleID,
			},
		}

		for _, c := range customerContracts {
			require.NoError(t, TestDb.CustomerContractStore.Create(c))
		}

		testCases := []struct {
			StartDate      int
			EndDate        int
			OptionParams   map[string]interface{}
			ExpectedLenCar int
		}{
			{StartDate: 0, EndDate: 2, OptionParams: map[string]interface{}{}, ExpectedLenCar: 2},
			{StartDate: 2, EndDate: 4, OptionParams: map[string]interface{}{}, ExpectedLenCar: 2},
			{StartDate: 10, EndDate: 12, OptionParams: map[string]interface{}{}, ExpectedLenCar: 3},
			{StartDate: 0, EndDate: 2, OptionParams: map[string]interface{}{"fuels": []string{string(model.FuelElectricity)}}, ExpectedLenCar: 1},
			{StartDate: 2, EndDate: 4, OptionParams: map[string]interface{}{"motions": []string{string(model.MotionAutomaticTransmission)}}, ExpectedLenCar: 1},
			{StartDate: 10, EndDate: 12, OptionParams: map[string]interface{}{"parking_lots": []string{string(model.ParkingLotGarage)}}, ExpectedLenCar: 2},
			{StartDate: 10, EndDate: 12, OptionParams: map[string]interface{}{"parking_lots": []string{string(model.ParkingLotGarage)}, "motions": []string{string(model.MotionManualTransmission)}, "number_of_seats": []int{15}}, ExpectedLenCar: 1},
			{StartDate: 10, EndDate: 12, OptionParams: map[string]interface{}{"parking_lots": []string{string(model.ParkingLotGarage)}, "motions": []string{string(model.MotionAutomaticTransmission)}}, ExpectedLenCar: 1},
			{StartDate: 0, EndDate: 24, OptionParams: map[string]interface{}{"parking_lots": []string{string(model.ParkingLotGarage)}, "motions": []string{string(model.MotionAutomaticTransmission)}}, ExpectedLenCar: 0},
		}

		toTime := func(hour int) time.Time {
			return now.Add(10 * time.Hour * time.Duration(hour))
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%d -> %d", tc.StartDate, tc.EndDate), func(t *testing.T) {
				foundCars, err := TestDb.CarStore.FindCars(toTime(tc.StartDate), toTime(tc.EndDate), tc.OptionParams)
				require.NoError(t, err)
				require.Len(t, foundCars, tc.ExpectedLenCar)
			})
		}
	})
}
