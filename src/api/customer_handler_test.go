package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
	cars := []*model.Car{
		{PartnerID: partner.ID, CarModelID: carModels[0].ID, Status: model.CarStatusActive},
		{PartnerID: partner.ID, CarModelID: carModels[1].ID, Status: model.CarStatusActive},
		{PartnerID: partner.ID, CarModelID: carModels[2].ID, Status: model.CarStatusActive},
		{PartnerID: partner.ID, CarModelID: carModels[1].ID, Status: model.CarStatusRejected},
		{PartnerID: partner.ID, CarModelID: carModels[2].ID, Status: model.CarStatusWaitingDelivery},
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

		query.Add("start_date", now.Format(time.RFC3339))
		query.Add("end_date", now.Add(time.Hour*time.Duration(3)).Format(time.RFC3339))
		req.URL.RawQuery = query.Encode()
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)

		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
		bz, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		var foundCars []*carResponse
		require.NoError(t, json.Unmarshal(bz, &foundCars))
		require.Len(t, foundCars, 3)
	})

	t.Run("start_date and two queries", func(t *testing.T) {
		req, err := http.NewRequest(route.Method, route.Path, nil)
		require.NoError(t, err)
		query := req.URL.Query()
		now := time.Now()

		query.Add("start_date", now.Format(time.RFC3339))
		query.Add("end_date", now.Add(time.Hour*time.Duration(3)).Format(time.RFC3339))
		query.Add("brand", "Lexus")
		req.URL.RawQuery = query.Encode()
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)

		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
		bz, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		var foundCars []*carResponse
		require.NoError(t, json.Unmarshal(bz, &foundCars))
		require.Len(t, foundCars, 1)
	})
}
