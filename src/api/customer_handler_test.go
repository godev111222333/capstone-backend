package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestCustomerHandler_FindCars(t *testing.T) {
	require.NoError(t, ResetDb(TestConfig.Database))

	carModels := []*model.CarModel{
		{Model: "Z9", Brand: "Lexus", NumberOfSeats: 15},
		{Model: "Z1000", Brand: "BMW", NumberOfSeats: 4},
		{Model: "Z800", Brand: "Kawasaki", NumberOfSeats: 7},
	}
	require.NoError(t, TestDb.CarModelStore.Create(carModels))
	partner := &model.Account{RoleID: model.RoleIDPartner, Email: "pn1@gmail@gmail.com", Status: model.AccountStatusActive}
	require.NoError(t, TestDb.AccountStore.Create(partner))
	PartnerContractRuleID := 1
	cars := []*model.Car{
		{PartnerID: partner.ID, CarModelID: carModels[0].ID, Status: model.CarStatusActive, LicensePlate: "232222", ParkingLot: model.ParkingLotHome, PartnerContractRuleID: PartnerContractRuleID},
		{PartnerID: partner.ID, CarModelID: carModels[1].ID, Status: model.CarStatusActive, LicensePlate: "242222", ParkingLot: model.ParkingLotGarage, PartnerContractRuleID: PartnerContractRuleID},
		{PartnerID: partner.ID, CarModelID: carModels[2].ID, Status: model.CarStatusActive, LicensePlate: "252222", ParkingLot: model.ParkingLotHome, PartnerContractRuleID: PartnerContractRuleID},
		{PartnerID: partner.ID, CarModelID: carModels[1].ID, Status: model.CarStatusRejected, LicensePlate: "262222", ParkingLot: model.ParkingLotGarage, PartnerContractRuleID: PartnerContractRuleID},
		{PartnerID: partner.ID, CarModelID: carModels[2].ID, Status: model.CarStatusWaitingDelivery, LicensePlate: "272222", ParkingLot: model.ParkingLotHome, PartnerContractRuleID: PartnerContractRuleID},
	}
	for _, c := range cars {
		require.NoError(t, TestDb.CarStore.Create(c))
	}
	route := TestServer.AllRoutes()[RouteCustomerFindCars]
	_, accessPayload := seedAccountAndLogin("cs1@gmail.com", "a", model.RoleIDCustomer)

	t.Run("only start_date and end_date", func(t *testing.T) {
		req, err := http.NewRequest(route.Method, route.Path, nil)
		require.NoError(t, err)
		query := req.URL.Query()
		now := time.Now()

		query.Add("start_date", now.Add(time.Second).Format(time.RFC3339))
		query.Add("end_date", now.AddDate(0, 0, 1).Add(time.Hour*time.Duration(3)).Format(time.RFC3339))
		req.URL.RawQuery = query.Encode()
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)

		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
		bz, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		var foundCars []*carResponse
		require.NoError(t, unmarshalFromCommResponse(bz, &foundCars))
		require.Len(t, foundCars, 3)
	})

	t.Run("start_date and two queries", func(t *testing.T) {
		req, err := http.NewRequest(route.Method, route.Path, nil)
		require.NoError(t, err)
		query := req.URL.Query()
		now := time.Now()

		query.Add("start_date", now.Add(time.Second).Format(time.RFC3339))
		query.Add("end_date", now.AddDate(0, 0, 1).Add(time.Hour*time.Duration(3)).Format(time.RFC3339))
		query.Add("brands", "Lexus")
		query.Add("number_of_seats", "4,7,15")
		req.URL.RawQuery = query.Encode()
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)

		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
		bz, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		var foundCars []*carResponse
		require.NoError(t, unmarshalFromCommResponse(bz, &foundCars))
		require.Len(t, foundCars, 1)
	})
}

func TestServer_HandleCustomerCalculateRentPricing(t *testing.T) {
	carModel := &model.CarModel{Brand: "xxx"}
	require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
	partner := &model.Account{FirstName: "pppppp", RoleID: model.RoleIDPartner, Status: model.AccountStatusActive, PhoneNumber: "9123192391239"}
	require.NoError(t, TestDb.AccountStore.Create(partner))
	car := &model.Car{CarModelID: carModel.ID, Price: 100_000, PartnerID: partner.ID, PartnerContractRuleID: 1}
	require.NoError(t, TestDb.CarStore.Create(car))

	route := TestServer.AllRoutes()[RouteCustomerCalculateRentingPrice]
	req, err := http.NewRequest(route.Method, route.Path, nil)
	require.NoError(t, err)
	query := req.URL.Query()
	now := time.Now()
	query.Add("car_id", strconv.Itoa(car.ID))
	query.Add("start_date", now.Format(time.RFC3339))
	query.Add("end_date", now.AddDate(0, 0, 3).Format(time.RFC3339))
	req.URL.RawQuery = query.Encode()
	_, accessPayload := seedAccountAndLogin("xxxxxxxx", "xxxx", model.RoleIDCustomer)
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)

	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	resp := &RentPricing{}
	bz, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)
	require.NoError(t, unmarshalFromCommResponse(bz, resp))

	require.Equal(t, 100_000, resp.RentPriceQuotation)
	require.Equal(t, 10_000, resp.InsurancePriceQuotation)
	require.Equal(t, 300_000, resp.TotalRentPriceAmount)
	require.Equal(t, 30_000, resp.TotalInsuranceAmount)
	require.Equal(t, 330_000, resp.TotalAmount)
	require.Equal(t, 99_000, resp.PrepaidAmount)
}

func TestServer_HandleCustomerGiveFeedback(t *testing.T) {
	carModel := &model.CarModel{Brand: "BMW"}
	require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
	partner, _ := seedAccountAndLogin("12231", "xxxx", model.RoleIDPartner)
	car := &model.Car{CarModelID: carModel.ID, LicensePlate: "kdjkas", PartnerID: partner.ID, PartnerContractRuleID: 1}
	require.NoError(t, TestDb.CarStore.Create(car))
	customer, authPayload := seedAccountAndLogin("2233", "xxxx", model.RoleIDCustomer)
	cusContract := &model.CustomerContract{
		CustomerID:             customer.ID,
		CarID:                  car.ID,
		Status:                 model.CustomerContractStatusCompleted,
		CustomerContractRuleID: 1,
	}
	require.NoError(t, TestDb.CustomerContractStore.Create(cusContract))

	route := TestServer.AllRoutes()[RouteCustomerGiveFeedback]
	reqBody := customerGiveFeedbackRequest{
		CustomerContractID: cusContract.ID,
		Content:            "good car",
		Rating:             5,
	}
	bz, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+authPayload.AccessToken)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	updateContract, err := TestDb.CustomerContractStore.FindByID(cusContract.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateContract)
	require.Equal(t, "good car", updateContract.FeedbackContent)
	require.Equal(t, 5, updateContract.FeedbackRating)
	require.Equal(t, model.FeedbackStatusActive, updateContract.FeedbackStatus)
}
