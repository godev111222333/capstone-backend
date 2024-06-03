package api

import (
	"bytes"
	"encoding/json"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPartnerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Register partner successfully", func(t *testing.T) {
		t.Parallel()

		route := TestServer.AllRoutes()[RouteRegisterPartner]
		body := `{
			"first_name": "Cuong",
			"last_name": "Nguyen Van",
			"phone_number": "8888",
			"email": "nguyenvancuong11@gmail.com",
			"identification_card_number": "6868",
			"password": "abcXYZ123"
		}`

		req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader([]byte(body)))
		require.NotNil(t, req)
		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestRegisterCar(t *testing.T) {
	t.Parallel()

	t.Run("register car successfully", func(t *testing.T) {
		t.Parallel()

		route := TestServer.AllRoutes()[RouteRegisterCar]
		carModel := &model.CarModel{
			Brand:         "Ferrari",
			Model:         "X9",
			Year:          2023,
			NumberOfSeats: 2,
			BasedPrice:    350_000,
		}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
		body := registerCarRequest{
			LicensePlate: "51A3",
			CarModelID:   carModel.ID,
			Motion:       model.MotionManualTransmission,
			Fuel:         model.FuelGas,
			ParkingLot:   model.ParkingLotHome,
			PeriodCode:   "1",
			Description:  "Super dude",
		}
		bz, _ := json.Marshal(body)
		req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))

		hashedPassword, _ := TestServer.hashVerifier.Hash("0000000")
		partner := &model.Account{
			RoleID:    model.RoleIDPartner,
			Email:     "bill@gmail.com",
			FirstName: "Bill Gate",
			Status:    model.AccountStatusActive,
			Password:  hashedPassword,
		}
		require.NoError(t, TestDb.AccountStore.Create(partner))
		accessToken := login(partner.Email, "0000000").AccessToken

		req.Header.Add(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken)

		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		bz, _ = io.ReadAll(recorder.Body)
		require.Equal(t, http.StatusOK, recorder.Code)
	})
}
