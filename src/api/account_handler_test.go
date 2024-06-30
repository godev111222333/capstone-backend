package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestAccountHandlerRawLogin(t *testing.T) {
	t.Run("HandleRawLogin", func(t *testing.T) {
		hashedPassword, err := TestServer.hashVerifier.Hash("password")
		require.NoError(t, err)

		acct := &model.Account{
			RoleID:      model.RoleIDPartner,
			FirstName:   "Cuong",
			LastName:    "Nguyen Van",
			Email:       "nguyenvancuong@gmail.com",
			Password:    hashedPassword,
			Status:      model.AccountStatusActive,
			PhoneNumber: "4324",
		}

		require.NoError(t, TestDb.AccountStore.Create(acct))

		route := TestServer.AllRoutes()[RouteRawLogin]
		body := `{
			"phone_number": "4324",
			"password": "password"
		}`
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader([]byte(body)))
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)

		bz, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		resp := &rawLoginResponse{}
		require.NoError(t, unmarshalFromCommResponse(bz, resp))

		require.NotEmpty(t, resp.AccessToken)
		require.NotEmpty(t, resp.AccessTokenExpiresAt)
		require.Equal(t, "Cuong", resp.User.FirstName)
		require.Equal(t, "Nguyen Van", resp.User.LastName)
		require.Equal(t, "nguyenvancuong@gmail.com", resp.User.Email)
		require.Equal(t, "partner", resp.User.Role)

		// Testing with required authorization route
		route = TestServer.AllRoutes()[RouteTestAuthorization]
		req, _ = http.NewRequest(route.Method, "/partner"+route.Path, nil)
		req.Header.Add(authorizationHeaderKey, authorizationTypeBearer+" "+resp.AccessToken)
		recorder = httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestRenewAccessToken(t *testing.T) {
	t.Run("renew access token successfully", func(t *testing.T) {
		// Create and login as partner role
		hashedPassword, err := TestServer.hashVerifier.Hash("password")
		require.NoError(t, err)

		acct := &model.Account{
			RoleID:                   model.RoleIDPartner,
			FirstName:                "ABCDE",
			LastName:                 "Nguyen Van",
			Email:                    "abcde@gmail.com",
			PhoneNumber:              "99999",
			Password:                 hashedPassword,
			Status:                   model.AccountStatusActive,
			IdentificationCardNumber: "asdasd",
		}

		require.NoError(t, TestDb.AccountStore.Create(acct))

		route := TestServer.AllRoutes()[RouteRawLogin]
		body := `{
			"phone_number": "99999",
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
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	t.Run("update profile successfully", func(t *testing.T) {
		t.Parallel()

		hashedPassword, err := TestServer.hashVerifier.Hash("3333")
		require.NoError(t, err)
		acct := &model.Account{
			RoleID:                   model.RoleIDPartner,
			FirstName:                "Tran Van",
			LastName:                 "Tuan",
			PhoneNumber:              "6754",
			Email:                    "1234@gmail.com",
			IdentificationCardNumber: "78910",
			Password:                 hashedPassword,
			DrivingLicense:           "8888",
			Status:                   model.AccountStatusActive,
		}
		require.NoError(t, TestDb.AccountStore.Create(acct))

		route := TestServer.AllRoutes()[RouteUpdateProfile]
		dob, err := time.Parse(time.RFC3339, "1998-06-20T20:30:00Z")
		require.NoError(t, err)
		accessToken := login(acct.PhoneNumber, "3333").AccessToken
		body := updateProfileRequest{
			FirstName:                "Son",
			LastName:                 "Le Thanh",
			Email:                    "new_email@gmail.com",
			DateOfBirth:              dob,
			IdentificationCardNumber: "111111111",
			DrivingLicense:           "111111111111",
			Password:                 "new_password",
		}
		bz, _ := json.Marshal(body)
		req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
		req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessToken)
		recorder := httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)

		updatedAcct, err := TestServer.store.AccountStore.GetByID(acct.ID)
		require.NoError(t, err)
		require.Equal(t, "Son", updatedAcct.FirstName)
		require.Equal(t, "Le Thanh", updatedAcct.LastName)
		require.Equal(t, "new_email@gmail.com", updatedAcct.Email)
		require.Equal(t, "111111111", updatedAcct.IdentificationCardNumber)
		require.Equal(t, "111111111111", updatedAcct.DrivingLicense)
		require.NoError(t, TestServer.hashVerifier.Compare(updatedAcct.Password, "new_password"))
	})
}

func TestUpdatePaymentInformation(t *testing.T) {
	t.Parallel()

	_, accessPayload := seedAccountAndLogin("sonle1@gmail.com", "9999", model.RoleIDPartner)

	route := TestServer.AllRoutes()[RouteUpdatePaymentInformation]
	reqBody := updatePaymentInfoRequest{
		BankNumber: "9999-9999-9999",
		BankOwner:  "Le Thanh Son",
		BankName:   "Ngan hang thuong mai co phan - VCB",
	}
	reqBz, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader(reqBz))
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)

	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	// get payment information after updating
	route = TestServer.AllRoutes()[RouteGetPaymentInformation]
	req, _ = http.NewRequest(route.Method, route.Path, nil)
	req.Header.Set(authorizationHeaderKey, authorizationTypeBearer+" "+accessPayload.AccessToken)
	recorder = httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	resp := &model.PaymentInformation{}
	bz, _ := io.ReadAll(recorder.Body)
	require.NoError(t, unmarshalFromCommResponse(bz, resp))

	require.Equal(t, "9999-9999-9999", resp.BankNumber)
	require.Equal(t, "Le Thanh Son", resp.BankOwner)
	require.Equal(t, "Ngan hang thuong mai co phan - VCB", resp.BankName)
}
