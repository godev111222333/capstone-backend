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
			RoleID:      model.RoleIDPartner,
			PhoneNumber: "bill@gmail.com",
			FirstName:   "Bill Gate",
			Status:      model.AccountStatusActive,
			Password:    hashedPassword,
		}
		require.NoError(t, TestDb.AccountStore.Create(partner))
		accessToken := login(partner.PhoneNumber, "0000000").AccessToken

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
	t.Run("update rental price", func(t *testing.T) {
		acct, accessPayload := seedAccountAndLogin("8989989", "xxx", model.RoleIDPartner)
		carModel := &model.CarModel{
			Brand: "Lambo",
		}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
		car := &model.Car{
			PartnerID:  acct.ID,
			CarModelID: carModel.ID,
			Price:      100_000,
			Status:     model.CarStatusPendingApplicationPendingPrice,
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
		Status:    model.PartnerContractStatusWaitingForAgreement,
	}
	require.NoError(t, TestDb.PartnerContractStore.Create(contract))

	route := TestServer.AllRoutes()[RoutePartnerAgreeContract]
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
	require.Equal(t, model.PartnerContractStatusAgreed, updatedContract.Status)
}

func TestRenderPartnerContract(t *testing.T) {
	t.Skip()

	previousPdfService := TestServer.pdfService
	TestServer.pdfService = NewPDFService(TestConfig.PDFService)
	defer func() {
		TestServer.pdfService = previousPdfService
	}()

	partner := &model.Account{
		FirstName:                "Cuong",
		LastName:                 "Tran Van",
		IdentificationCardNumber: "123456789",
		Email:                    "ppp1@gmail.com",
		Status:                   model.AccountStatusActive,
		RoleID:                   model.RoleIDPartner,
	}
	require.NoError(t, TestDb.AccountStore.Create(partner))
	carModel := &model.CarModel{Brand: "BMW", Model: "VIP 2", NumberOfSeats: 4, Year: 2024}
	require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
	car := &model.Car{
		PartnerID:    partner.ID,
		CarModelID:   carModel.ID,
		LicensePlate: "96A1",
		Price:        400_000,
		Period:       6,
	}
	require.NoError(t, TestDb.CarStore.Create(car))

	now := time.Now()
	contract := &model.PartnerContract{
		CarID:     car.ID,
		StartDate: now,
		EndDate:   now.AddDate(0, 6, 0),
	}
	require.NoError(t, TestDb.PartnerContractStore.Create(contract))
	ct, err := TestDb.PartnerContractStore.GetByCarID(car.ID)
	require.NoError(t, err)

	require.NoError(t, TestServer.RenderPartnerContractPDF(partner, &ct.Car))

	updatedContract, err := TestDb.PartnerContractStore.GetByCarID(car.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedContract.Url)
}

func TestRenderCustomerContract(t *testing.T) {
	t.Skip()

	previousPdfService := TestServer.pdfService
	TestServer.pdfService = NewPDFService(TestConfig.PDFService)
	defer func() {
		TestServer.pdfService = previousPdfService
	}()

	partner := &model.Account{
		FirstName:                "Khanh",
		LastName:                 "Tran Van",
		IdentificationCardNumber: "987678976",
		Email:                    "khanh@gmail.com",
		Status:                   model.AccountStatusActive,
		RoleID:                   model.RoleIDPartner,
	}
	require.NoError(t, TestDb.AccountStore.Create(partner))
	carModel := &model.CarModel{Brand: "BMW", Model: "VIP 2", NumberOfSeats: 4, Year: 2024}
	require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
	car := &model.Car{
		PartnerID:    partner.ID,
		CarModelID:   carModel.ID,
		LicensePlate: "96A1",
		Price:        400_000,
		Period:       6,
	}
	require.NoError(t, TestDb.CarStore.Create(car))
	var err error
	car, err = TestDb.CarStore.GetByID(car.ID)
	require.NoError(t, err)

	customer := &model.Account{
		FirstName:                "Toan",
		LastName:                 "Le Thanh",
		IdentificationCardNumber: "88888888",
		Email:                    "toan@gmail.com",
		Status:                   model.AccountStatusActive,
		RoleID:                   model.RoleIDCustomer,
	}
	require.NoError(t, TestDb.AccountStore.Create(customer))

	now := time.Now()
	contract := &model.CustomerContract{
		CustomerID:              customer.ID,
		CarID:                   car.ID,
		RentPrice:               car.Price * 3,
		StartDate:               now,
		EndDate:                 now.AddDate(0, 0, 3),
		Status:                  model.CustomerContractStatusWaitingContractAgreement,
		InsuranceAmount:         car.Price / 10,
		CollateralType:          model.CollateralTypeCash,
		IsReturnCollateralAsset: false,
	}
	require.NoError(t, TestDb.CustomerContractStore.Create(contract))

	require.NoError(t, TestServer.RenderCustomerContractPDF(customer, car, contract))
	updatedContract, err := TestDb.CustomerContractStore.FindByID(contract.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedContract.Url)
}
