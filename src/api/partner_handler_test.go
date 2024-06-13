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
	"time"
)

func TestRegisterPartnerHandler(t *testing.T) {
	t.Skip()

	t.Run("Register partner successfully", func(t *testing.T) {
		t.Parallel()

		route := TestServer.AllRoutes()[RouteRegisterPartner]
		body := `{
			"first_name": "Cuong",
			"last_name": "Nguyen Van",
			"phone_number": "8888",
			"email": "nguyenvancuong11@gmail.com",
			"password": "abcXYZ123"
		}`

		req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader([]byte(body)))
		require.NotNil(t, req)
		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestRegisterCarHandler(t *testing.T) {
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

		// get registered cars
		recorder = httptest.NewRecorder()
		route = TestServer.AllRoutes()[RouteGetRegisteredCars]
		req, _ = http.NewRequest(route.Method, route.Path, nil)
		q := req.URL.Query()
		q.Set("offset", "0")
		q.Set("limit", "1")
		q.Set("car_status", string(model.CarStatusPendingApplicationPendingCarImages))
		req.URL.RawQuery = q.Encode()
		req.Header.Add(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken)
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)

		resp := &getRegisteredCarsResponse{}
		bz, _ = io.ReadAll(recorder.Body)
		require.NoError(t, json.Unmarshal(bz, resp))
		require.Len(t, resp.Cars, 1)
	})
}

func TestUpdateRentalPriceHandler(t *testing.T) {
	t.Parallel()

	t.Run("update rental price", func(t *testing.T) {
		t.Parallel()

		acct, accessPayload := seedAccountAndLogin("minh@gmail.com", "xxx", model.RoleIDPartner)
		carModel := &model.CarModel{
			Brand: "Lambo",
		}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
		car := &model.Car{
			PartnerID:  acct.ID,
			CarModelID: carModel.ID,
			Price:      100_000,
			Status:     model.CarStatusPendingApproval,
		}
		require.NoError(t, TestDb.CarStore.Create(car))

		route := TestServer.AllRoutes()[RouteUpdateRentalPrice]
		body := updateRentalPriceRequest{
			CarID:    car.ID,
			NewPrice: 200_000,
		}
		bz, _ := json.Marshal(body)
		req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)
		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)

		updatedCar, err := TestServer.store.CarStore.GetByID(car.ID)
		require.NoError(t, err)
		require.Equal(t, 200_000, updatedCar.Price)
	})
}

func TestSignContract(t *testing.T) {
	t.Parallel()

	carModel := &model.CarModel{Brand: "Ok"}
	require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
	partner, accessPayload := seedAccountAndLogin("partner1", "aa", model.RoleIDPartner)
	car := &model.Car{
		PartnerID:    partner.ID,
		CarModelID:   carModel.ID,
		LicensePlate: "89A8",
	}
	require.NoError(t, TestDb.CarStore.Create(car))
	period := 3
	contract := &model.PartnerContract{
		CarID:     car.ID,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, period, 0),
		Url:       FakePDF,
		Status:    model.PartnerContractStatusWaitingForSigning,
	}
	require.NoError(t, TestDb.PartnerContractStore.Create(contract))

	route := TestServer.AllRoutes()[RoutePartnerSignContract]
	r := partnerSignContractRequest{CarID: car.ID}
	bz, err := json.Marshal(r)
	require.NoError(t, err)
	req, err := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	updatedContract, err := TestDb.PartnerContractStore.GetByCarID(car.ID)
	require.NoError(t, err)
	require.Equal(t, model.PartnerContractStatusSigned, updatedContract.Status)
}
