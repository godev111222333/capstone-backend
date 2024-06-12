package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/godev111222333/capstone-backend/src/model"
)

func seedAccountAndLogin(email, password string, role model.RoleID) (*model.Account, *rawLoginResponse) {
	h, _ := TestServer.hashVerifier.Hash(password)
	acct := &model.Account{
		Email:    email,
		Password: h,
		RoleID:   role,
		Status:   model.AccountStatusActive,
	}
	_ = TestDb.AccountStore.Create(acct)

	return acct, login(email, password)
}

func loginAdmin() *rawLoginResponse {
	return login("admin", "admin")
}

func login(email, password string) *rawLoginResponse {
	route := TestServer.AllRoutes()[RouteRawLogin]
	body := rawLoginRequest{
		Email:    email,
		Password: password,
	}
	bz, _ := json.Marshal(body)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
	TestServer.route.ServeHTTP(recorder, req)
	bz, _ = io.ReadAll(recorder.Body)
	res := &rawLoginResponse{}
	_ = json.Unmarshal(bz, res)
	return res
}
