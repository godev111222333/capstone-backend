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

		resp := struct {
			Models  []model.CarModel `json:"models"`
			Periods []int            `json:"periods"`
		}{}
		bz, _ := io.ReadAll(recorder.Body)
		require.NoError(t, json.Unmarshal(bz, &resp))
		require.NotEmpty(t, resp.Models)
		require.Equal(t, resp.Periods, []int{1, 3, 6, 12})
	})
}
