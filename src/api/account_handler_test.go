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
			Password:                 hashedPassword,
			Status:                   model.AccountStatusActive,
			IdentificationCardNumber: "asdasd",
		}

		require.NoError(t, TestDb.AccountStore.Create(acct))

		route := TestServer.AllRoutes()[RouteRawLogin]
		body := `{
			"email": "abcde@gmail.com",
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
		require.NotEmpty(t, resp.RefreshToken)

		refreshTokenPayload, err := TestServer.tokenMaker.VerifyToken(resp.RefreshToken)
		require.NoError(t, err)
		session, err := TestDb.SessionStore.GetSession(refreshTokenPayload.ID)
		require.NoError(t, err)
		require.NotEmpty(t, session)

		route = TestServer.AllRoutes()[RouteRenewAccessToken]
		renewBody := renewAccessTokenRequest{RefreshToken: resp.RefreshToken}
		bz, err = json.Marshal(renewBody)
		require.NoError(t, err)
		req, err = http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
		require.NoError(t, err)
		recorder = httptest.NewRecorder()
		TestServer.route.ServeHTTP(recorder, req)
		bz, err = io.ReadAll(recorder.Body)
		require.NoError(t, err)

		renewResp := renewAccessTokenResponse{}
		require.NoError(t, json.Unmarshal(bz, &renewResp))
		require.NotEmpty(t, renewResp.AccessToken)
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
		accessToken := login(acct.Email, "3333").AccessToken
		body := updateProfileRequest{
			FirstName:                "Son",
			LastName:                 "Le Thanh",
			PhoneNumber:              "0123456",
			DateOfBirth:              dob,
			IdentificationCardNumber: "1111",
			DrivingLicense:           "2222",
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
		require.Equal(t, "0123456", updatedAcct.PhoneNumber)
		require.Equal(t, "1111", updatedAcct.IdentificationCardNumber)
		require.Equal(t, "2222", updatedAcct.DrivingLicense)
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
	require.NoError(t, json.Unmarshal(bz, resp))

	require.Equal(t, "9999-9999-9999", resp.BankNumber)
	require.Equal(t, "Le Thanh Son", resp.BankOwner)
	require.Equal(t, "Ngan hang thuong mai co phan - VCB", resp.BankName)
}

func TestAccountStore_Get(t *testing.T) {
	accounts := []*model.Account{
		{Email: "acc1@gmail.com", FirstName: "First 1", LastName: "Last 1", PhoneNumber: "phone1", RoleID: model.RoleIDCustomer, Status: model.AccountStatusActive},
		{Email: "acc2@gmail.com", FirstName: "First 2", LastName: "Last 2", PhoneNumber: "phone2", RoleID: model.RoleIDPartner, Status: model.AccountStatusActive},
		{Email: "acc3@gmail.com", FirstName: "First 2", LastName: "Last 2", PhoneNumber: "phone3", RoleID: model.RoleIDCustomer, Status: model.AccountStatusActive},
	}

	for _, acct := range accounts {
		require.NoError(t, TestDb.AccountStore.Create(acct))
	}

	t.Run("list all active accounts", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 3)
	})

	t.Run("list all partners", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, model.RoleNamePartner, "", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
	})

	t.Run("list all customers", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, model.RoleNameCustomer, "", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 2)
	})

	t.Run("get account 1", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "phone1", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
	})

	t.Run("get by Last", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "Last", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 3)
	})

	t.Run("get by phone1", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "phone1", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
	})

	t.Run("get by email", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "acc3@gmail.com", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
	})
}
