package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func getWirelessInterfaceRpi() (string, error) {
	return "wlan0", nil
}

func newTestAppHandler() testAppHandler {
	return testAppHandler{
		NetworkSettings: nil,
		H:               nil,
		SetAPModeCalls:  []bool{},
	}
}

type testAppHandler struct {
	NetworkSettings *models.NetworkSettings
	H               func(IAppHandler, http.ResponseWriter, *http.Request) (int, error)
	SetAPModeCalls  []bool
}

func (ah testAppHandler) GetContext() *models.AppContext {
	v := time.Time{}
	return &models.AppContext{
		Time: &v,
	}
}

func (ah testAppHandler) WriteConfig() {}

func (ah *testAppHandler) SetAPMode(v bool) error {
	ah.SetAPModeCalls = append(ah.SetAPModeCalls, v)
	return nil
}

func (ah testAppHandler) GetNetworkSettings() *models.NetworkSettings {
	return ah.NetworkSettings
}
func (ah testAppHandler) GetTime() *time.Time {
	v := time.Now()
	return &v
}

func (ah *testAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := ah.H(ah, w, r)
	handleHandlerResult(status, err, w, r)
}

func doRequest(
	testAppHandler *testAppHandler,
	expectedStatus int,
	t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	testAppHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	return rr
}
