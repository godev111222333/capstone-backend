package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/godev111222333/capstone-backend/src/model"
)

func seedAccountAndLogin(phoneNumber, password string, role model.RoleID) (*model.Account, *rawLoginResponse) {
	h, _ := TestServer.hashVerifier.Hash(password)
	acct := &model.Account{
		PhoneNumber: phoneNumber,
		Password:    h,
		RoleID:      role,
		Status:      model.AccountStatusActive,
	}
	_ = TestDb.AccountStore.Create(acct)

	return acct, login(phoneNumber, password)
}

func loginAdmin() *rawLoginResponse {
	return login("admin", "admin")
}

func login(phoneNumber, password string) *rawLoginResponse {
	route := TestServer.AllRoutes()[RouteRawLogin]
	body := rawLoginRequest{
		PhoneNumber: phoneNumber,
		Password:    password,
	}
	bz, _ := json.Marshal(body)
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest(route.Method, route.Path, bytes.NewReader(bz))
	TestServer.route.ServeHTTP(recorder, req)
	bz, _ = io.ReadAll(recorder.Body)
	res := &rawLoginResponse{}
	_ = unmarshalFromCommResponse(bz, res)
	return res
}

func unmarshalFromCommResponse(respBody []byte, data any) error {
	commResponse := &CommResponse{}
	if err := json.Unmarshal(respBody, commResponse); err != nil {
		return err
	}
	bz, err := json.Marshal(commResponse.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bz, data)
}
