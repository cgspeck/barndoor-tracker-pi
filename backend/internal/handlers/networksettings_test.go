package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"

	"github.com/pyros2097/cupaloy"
)

func TestNetworkSettingsHandler(t *testing.T) {
	rr := doRequest(NetworkSettingsHandler, getWirelessInterfaceRpi, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}

func doPost(
	appHandlerFunc func(IAppHandler, http.ResponseWriter, *http.Request) (int, error),
	body string,
	appContext *models.AppContext,
	t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := AppHandler{appContext, appHandlerFunc}

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	return rr
}
func TestNetworkSettingsHandlerPost(t *testing.T) {
	body := `{
	"accessPointMode": true
}
`
	rr := doPost(NetworkSettingsHandler, body, &models.AppContext{}, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}
