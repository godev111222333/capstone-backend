package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

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

	route := TestServer.AllRoutes()[RouteAdminGetCarDetails]
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
	t.Parallel()

	t.Run("approve successfully", func(t *testing.T) {
		t.Parallel()

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

		route := TestServer.AllRoutes()[RouteAdminApproveCar]
		req, err := http.NewRequest(
			route.Method,
			strings.Replace(route.Path, ":id", strconv.Itoa(car.ID), 1),
			nil,
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
