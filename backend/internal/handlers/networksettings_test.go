package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

type networkSettingsTestAppHandler struct {
	NetworkSettings *models.NetworkSettings
	H               func(IAppHandler, http.ResponseWriter, *http.Request) (int, error)
	SetAPModeCalls  []bool
}

func (ah networkSettingsTestAppHandler) GetContext() *models.AppContext {
	return &models.AppContext{}
}

func (ah networkSettingsTestAppHandler) WriteConfig() {}

func (ah *networkSettingsTestAppHandler) SetAPMode(v bool) {
	ah.SetAPModeCalls = append(ah.SetAPModeCalls, v)
}
func (ah networkSettingsTestAppHandler) GetNetworkSettings() *models.NetworkSettings {
	return ah.NetworkSettings
}
func (ah networkSettingsTestAppHandler) GetTime() *time.Time {
	v := time.Now()
	return &v
}

func doNetworkSettingsPost(
	body string,
	networkSettingsTestAppHandler *networkSettingsTestAppHandler,
	t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	networkSettingsTestAppHandler.H(networkSettingsTestAppHandler, rr, req)

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

	handler := networkSettingsTestAppHandler{&models.NetworkSettings{
		AccessPointMode: false,
	}, NetworkSettingsHandler, []bool{}}

	rr := doNetworkSettingsPost(body, &handler, t)

	// Check the response body is what we expect.
	err := cupaloy.Snapshot(rr)
	if err != nil {
		t.Error(err)
	}

	if len(handler.SetAPModeCalls) != 1 {
		t.Errorf("Expected call to SetAPMode")
	}

	if handler.SetAPModeCalls[0] != true {
		t.Errorf("Expected true call to SetAPMode")
	}
}
