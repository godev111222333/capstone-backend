package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
