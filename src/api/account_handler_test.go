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

func TestAccountHandlerRawLogin(t *testing.T) {
	t.Parallel()

	t.Run("HandleRawLogin", func(t *testing.T) {
		hashedPassword, err := TestServer.hashVerifier.Hash("password")
		require.NoError(t, err)

		acct := &model.Account{
			RoleID:    model.RoleIDPartner,
			FirstName: "Cuong",
			LastName:  "Nguyen Van",
			Email:     "nguyenvancuong@gmail.com",
			Password:  hashedPassword,
			Status:    model.AccountStatusEnable,
		}

		require.NoError(t, TestDb.AccountStore.Create(acct))

		route := TestServer.AllRoutes()[RouteRawLogin]
		body := `{
			"email": "nguyenvancuong@gmail.com",
			"password": "password"
		}`
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader([]byte(body)))
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)

		bz, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		resp := rawLoginResponse{}
		require.NoError(t, json.Unmarshal(bz, &resp))

		require.NotEmpty(t, resp.AccessToken)
		require.NotEmpty(t, resp.AccessTokenExpiresAt)
		require.NotEmpty(t, resp.RefreshToken)
		require.NotEmpty(t, resp.RefreshTokenExpiresAt)
		require.Equal(t, "Cuong", resp.User.FirstName)
		require.Equal(t, "Nguyen Van", resp.User.LastName)
		require.Equal(t, "nguyenvancuong@gmail.com", resp.User.Email)
		require.Equal(t, "partner", resp.User.Role)

		// Testing with required authorization route
		route = TestServer.AllRoutes()[RouteTestAuthorization]
		req, _ = http.NewRequest(route.Method, route.Path, nil)
		req.Header.Add(authorizationHeaderKey, authorizationTypeBearer+" "+resp.AccessToken)
		recorder = httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
	})
}
