package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataHandler(t *testing.T) {
	t.Parallel()

	route := TestServer.AllRoutes()[RouteGetBankMetadata]
	req, _ := http.NewRequest(route.Method, route.Path, nil)

	recorder := httptest.NewRecorder()
	TestServer.route.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	resp := struct {
		Banks []string `json:"banks"`
	}{}
	bz, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bz, &resp))
	require.Len(t, resp.Banks, 97)
}
