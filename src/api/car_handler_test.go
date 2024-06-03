package api

import (
	"encoding/json"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCarHandler(t *testing.T) {
	t.Parallel()

	t.Run("get all car models", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{
			{
				Brand:         "Audi",
				Model:         "A8",
				Year:          2024,
				NumberOfSeats: 2,
				BasedPrice:    500_000,
			},
		}))

		route := TestServer.AllRoutes()[RouteGetRegisterCarMetadata]
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest(route.Method, route.Path, nil)
		TestServer.route.ServeHTTP(recorder, req)

		resp := registerCarMetadataResponse{}
		bz, _ := io.ReadAll(recorder.Body)
		require.NoError(t, json.Unmarshal(bz, &resp))
		require.NotEmpty(t, resp.Models)
		require.Equal(t, []OptionResponse{
			{Code: "1", Text: "1 tháng"},
			{Code: "3", Text: "3 tháng"},
			{Code: "6", Text: "6 tháng"},
			{Code: "12", Text: "12 tháng"},
		}, resp.Periods)
		require.Equal(t, []OptionResponse{
			{Code: "gas", Text: "Xăng"},
			{Code: "oil", Text: "Dầu"},
			{Code: "electricity", Text: "Điện"},
		}, resp.Fuels)
		require.Equal(t, []OptionResponse{
			{Code: "automatic_transmission", Text: "Số tự động"},
			{Code: "manual_transmission", Text: "Số sàn"},
		}, resp.Motions)
	})
}
