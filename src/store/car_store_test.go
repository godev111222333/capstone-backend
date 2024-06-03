package store

import (
	"testing"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
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
}
