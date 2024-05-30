package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPartnerHandler(t *testing.T) {
	t.Parallel()

	t.Run("Create partner successfully", func(t *testing.T) {
		t.Parallel()

		route := TestServer.AllRoutes()[RouteRegisterPartner]
		body := `{
			"first_name": "Cuong",
			"last_name": "Nguyen Van",
			"phone_number": "8888",
			"email": "nguyenvancuongxyz@gmail.com",
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
