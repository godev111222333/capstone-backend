package store

import (
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestCarStore(t *testing.T) {
	t.Parallel()

	t.Run("create car successfully", func(t *testing.T) {
		t.Parallel()

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
			PartnerID:    partner.ID,
			CarModelID:   carModel.ID,
			LicensePlate: "7777",
			ParkingLot:   model.ParkingLotHome,
			Description:  "Beautiful car",
			Fuel:         model.FuelElectricity,
			Motion:       model.MotionAutomaticTransmission,
			Price:        550_000,
			Status:       model.CarStatusActive,
		}
		require.NoError(t, TestDb.CarStore.Create(car))
	})

	t.Run("get owned car successfully", func(t *testing.T) {
		t.Parallel()

		partner := &model.Account{
			RoleID:    model.RoleIDPartner,
			Email:     "cuongdola@gmail.com",
			FirstName: "Cuong dola",
			Status:    model.AccountStatusActive,
		}
		require.NoError(t, TestDb.AccountStore.Create(partner))
		carModel := &model.CarModel{
			Brand: "Bugatti",
		}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))

		for i := 1; i <= 2; i++ {
			car := &model.Car{
				PartnerID:    partner.ID,
				CarModelID:   carModel.ID,
				LicensePlate: "86A" + strconv.Itoa(i),
				Status:       model.CarStatusActive,
			}
			require.NoError(t, TestDb.CarStore.Create(car))
		}
		car := &model.Car{
			PartnerID:    partner.ID,
			CarModelID:   carModel.ID,
			LicensePlate: "86AX",
			Status:       model.CarStatusPendingApproval,
		}
		require.NoError(t, TestDb.CarStore.Create(car))

		cars, err := TestDb.CarStore.GetByPartner(partner.ID, 0, 0, model.CarStatusNoFilter)
		require.NoError(t, err)
		require.Len(t, cars, 3)
		cars, err = TestDb.CarStore.GetByPartner(partner.ID, 0, 2, model.CarStatusPendingApproval)
		require.NoError(t, err)
		require.Len(t, cars, 1)
	})
}
