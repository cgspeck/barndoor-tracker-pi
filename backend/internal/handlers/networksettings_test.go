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
	handler := newTestAppHandler()
	handler.H = NetworkSettingsHandler

	handler.NetworkSettings = &models.NetworkSettings{
		AccessPointMode: true,
		APSettings: &models.APSettings{
			Channel: 11,
			Key:     "",
			SSID:    "barndoor-tracker",
		},
		AvailableNetworks: make([]*models.AvailableNetwork, 0),
		ManagementEnabled: false,
		WirelessInterface: "wlan0",
		WirelessProfiles:  make([]*models.WirelessProfile, 0),
	}

	rr := doRequest(&handler, 200, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}

func doNetworkSettingsPost(
	body string,
	testAppHandler *testAppHandler,
	expectedStatus int,
	t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	testAppHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}
	return rr
}
func TestNetworkSettingsHandlerPost(t *testing.T) {
	body := `{
	"accessPointMode": true
}
`

	handler := newTestAppHandler()
	handler.H = NetworkSettingsHandler
	handler.NetworkSettings = &models.NetworkSettings{
		AccessPointMode:   false,
		ManagementEnabled: true,
	}

	rr := doNetworkSettingsPost(body, &handler, http.StatusOK, t)

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

func TestNetworkSettingsHandlerPostManagementDisabled(t *testing.T) {
	body := `{
	"accessPointMode": true
}
`

	handler := newTestAppHandler()
	handler.H = NetworkSettingsHandler
	handler.NetworkSettings = &models.NetworkSettings{
		AccessPointMode:   false,
		ManagementEnabled: false,
	}

	rr := doNetworkSettingsPost(body, &handler, http.StatusBadRequest, t)

	// Check the response body is what we expect.
	err := cupaloy.Snapshot(rr)
	if err != nil {
		t.Error(err)
	}

	if len(handler.SetAPModeCalls) != 0 {
		t.Errorf("Expected no call to SetAPMode")
	}
}
