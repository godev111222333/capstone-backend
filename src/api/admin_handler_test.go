package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestAdminHandler_GarageConfigs(t *testing.T) {
	t.Parallel()

	route := TestServer.AllRoutes()[RouteUpdateGarageConfigs]
	_, accessPayload := seedAccountAndLogin("admin1", "admin1", model.RoleIDAdmin)

	r := updateGarageConfigRequest{
		Max4Seats:  3,
		Max7Seats:  6,
		Max15Seats: 9,
	}
	bz, _ := json.Marshal(r)
	req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)
	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	route = TestServer.AllRoutes()[RouteGetGarageConfigs]
	req, _ = http.NewRequest(route.Method, route.Path, nil)
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)
	recorder = httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	bz, _ = io.ReadAll(recorder.Body)

	resp := getGarageConfigResponse{}
	require.NoError(t, json.Unmarshal(bz, &resp))

	require.Equal(t, 3, resp.Max4Seats)
	require.Equal(t, 6, resp.Max7Seats)
	require.Equal(t, 9, resp.Max15Seats)
	require.Equal(t, 18, resp.Total)
}

func TestAdminHandler_GetCar(t *testing.T) {
	t.Parallel()

	carModel := &model.CarModel{
		Brand: "Toyota",
	}
	require.NoError(t, TestServer.store.CarModelStore.Create([]*model.CarModel{carModel}))
	partner, _ := seedAccountAndLogin("parter@gmail.com", "aa", model.RoleIDPartner)
	car := &model.Car{
		PartnerID:    partner.ID,
		CarModelID:   carModel.ID,
		LicensePlate: "59A33",
		Status:       model.CarStatusActive,
	}
	require.NoError(t, TestDb.CarStore.Create(car))

	adminAuthPayload := login("admin", "admin")

	route := TestServer.AllRoutes()[RouteGetCarDetail]
	req, err := http.NewRequest(
		route.Method,
		strings.Replace(route.Path, ":id", strconv.Itoa(car.ID), 1),
		nil,
	)
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+adminAuthPayload.AccessToken)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	car = &model.Car{}
	bz, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bz, car))
	require.Equal(t, "59A33", car.LicensePlate)
}

func TestHandleApproveCar(t *testing.T) {
	t.Run("reject car", func(t *testing.T) {
		carModel := &model.CarModel{Brand: "toyota"}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
		partner, _ := seedAccountAndLogin("partner_vip", "aaa", model.RoleIDPartner)
		car := &model.Car{
			PartnerID:    partner.ID,
			CarModelID:   carModel.ID,
			LicensePlate: "69A1",
			Status:       model.CarStatusPendingApproval,
		}
		require.NoError(t, TestDb.CarStore.Create(car))
		accessToken := loginAdmin()

		reqB := adminApproveOrRejectRequest{CarID: car.ID, Action: "reject"}
		reqBz, err := json.Marshal(reqB)
		require.NoError(t, err)
		route := TestServer.AllRoutes()[RouteAdminApproveCar]
		req, err := http.NewRequest(
			route.Method,
			route.Path,
			bytes.NewReader(reqBz),
		)
		require.NoError(t, err)
		fmt.Println(accessToken)
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken.AccessToken)
		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		bz, _ := io.ReadAll(recorder.Body)
		fmt.Println(string(bz))
		require.Equal(t, http.StatusOK, recorder.Code)
		updatedCar, err := TestDb.CarStore.GetByID(car.ID)
		require.NoError(t, err)
		require.Equal(t, model.CarStatusRejected, updatedCar.Status)
	})

	t.Run("approve registration car", func(t *testing.T) {
		previousPdfService := TestServer.pdfService
		TestServer.pdfService = &MockPDFService{}
		defer func() {
			TestServer.pdfService = previousPdfService
		}()

		carModel := &model.CarModel{Brand: "toyota"}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
		partner, _ := seedAccountAndLogin("partner_vip2", "aaa", model.RoleIDPartner)
		car := &model.Car{
			PartnerID:    partner.ID,
			CarModelID:   carModel.ID,
			LicensePlate: "69A12",
			Status:       model.CarStatusPendingApproval,
		}
		require.NoError(t, TestDb.CarStore.Create(car))
		accessToken := loginAdmin()

		reqB := adminApproveOrRejectRequest{CarID: car.ID, Action: "approve_register"}
		reqBz, err := json.Marshal(reqB)
		require.NoError(t, err)
		route := TestServer.AllRoutes()[RouteAdminApproveCar]
		req, err := http.NewRequest(
			route.Method,
			route.Path,
			bytes.NewReader(reqBz),
		)
		require.NoError(t, err)
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken.AccessToken)
		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
		updatedCar, err := TestDb.CarStore.GetByID(car.ID)
		require.NoError(t, err)
		require.Equal(t, model.CarStatusApproved, updatedCar.Status)
		time.Sleep(2 * time.Second)
	})

	t.Run("approve delivery car", func(t *testing.T) {
		carModel := &model.CarModel{Brand: "toyota"}
		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
		partner, _ := seedAccountAndLogin("partner_vip3", "aaa", model.RoleIDPartner)
		car := &model.Car{
			PartnerID:    partner.ID,
			CarModelID:   carModel.ID,
			LicensePlate: "69A11111",
			Status:       model.CarStatusWaitingDelivery,
		}
		require.NoError(t, TestDb.CarStore.Create(car))
		contract := &model.PartnerContract{
			CarID:     car.ID,
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 3, 0),
			Status:    model.PartnerContractStatusAgreed,
		}
		require.NoError(t, TestDb.PartnerContractStore.Create(contract))
		accessToken := loginAdmin()

		reqB := adminApproveOrRejectRequest{CarID: car.ID, Action: "approve_delivery"}
		reqBz, err := json.Marshal(reqB)
		require.NoError(t, err)
		route := TestServer.AllRoutes()[RouteAdminApproveCar]
		req, err := http.NewRequest(
			route.Method,
			route.Path,
			bytes.NewReader(reqBz),
		)
		require.NoError(t, err)
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken.AccessToken)
		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
		updatedCar, err := TestDb.CarStore.GetByID(car.ID)
		require.NoError(t, err)
		require.Equal(t, model.CarStatusActive, updatedCar.Status)
	})
}

func TestAdminHandler_CustomerPayments(t *testing.T) {
	carModel := &model.CarModel{Brand: "Ok"}
	require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
	partner, _ := seedAccountAndLogin("22233", "2222", model.RoleIDPartner)
	car := &model.Car{CarModelID: carModel.ID, PartnerID: partner.ID, Status: model.CarStatusActive, LicensePlate: "2222"}
	require.NoError(t, TestDb.CarStore.Create(car))
	customer, _ := seedAccountAndLogin("22234", "2222", model.RoleIDCustomer)
	customerContract := &model.CustomerContract{CarID: car.ID, CustomerID: customer.ID}
	require.NoError(t, TestDb.CustomerContractStore.Create(customerContract))
	route := TestServer.AllRoutes()[RouteAdminCreateCustomerPayment]

	accessPayload := loginAdmin()
	requestBody := adminCreateCustomerPaymentRequest{
		CustomerContractID: customerContract.ID,
		PaymentType:        model.PaymentTypeRemainingPay,
		Amount:             200_000,
		Note:               "vip",
	}
	bz, err := json.Marshal(requestBody)
	require.NoError(t, err)
	req, err := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
	req.Header.Add(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	route = TestServer.AllRoutes()[RouteAdminGetCustomerPayments]
	req, err = http.NewRequest(route.Method, route.Path, nil)
	req.Header.Add(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)
	query := req.URL.Query()
	query.Set("customer_contract_id", strconv.Itoa(customerContract.ID))
	query.Set("payment_status", string(model.PaymentStatusPending))
	req.URL.RawQuery = query.Encode()

	recorder = httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	bz, err = io.ReadAll(recorder.Body)
	payments := make([]*customerPaymentResponse, 0)
	require.NoError(t, json.Unmarshal(bz, &payments))
	require.Len(t, payments, 1)
}
