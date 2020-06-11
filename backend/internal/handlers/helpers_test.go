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
	lsc := make([]map[string]interface{}, 0)
	asc := make([]map[string]interface{}, 0)

	return testAppHandler{
		AlignStatus:      nil,
		TrackStatus:      nil,
		APSettings:       nil,
		NetworkSettings:  nil,
		LocationSettings: nil,
		Flags:            nil,
		H:                nil,
		// the following are so the tests can inspect what the handlers do
		SetAPModeCalls:           []bool{},
		SetLocationSettingsCalls: lsc,
		SetAPSettingsCalls:       asc,
	}
}

type testAppHandler struct {
	APSettings               *models.APSettings
	AlignStatus              *models.AlignStatus
	TrackStatus              *models.TrackStatus
	NetworkSettings          *models.NetworkSettings
	LocationSettings         *models.LocationSettings
	Flags                    *models.Flags
	H                        func(IAppHandler, http.ResponseWriter, *http.Request) (int, error)
	SetAPModeCalls           []bool
	SetLocationSettingsCalls []map[string]interface{}
	SetAPSettingsCalls       []map[string]interface{}
}

func (ah testAppHandler) GetContext() *models.AppContext {
	v := time.Time{}
	return &models.AppContext{
		Time:            &v,
		AlignStatus:     ah.AlignStatus,
		Flags:           ah.Flags,
		NetworkSettings: ah.NetworkSettings,
	}
}

func (ah *testAppHandler) SetAPMode(v bool) error {
	ah.SetAPModeCalls = append(ah.SetAPModeCalls, v)
	return nil
}

func (ah *testAppHandler) SetLocationSettings(input map[string]interface{}) error {
	ah.SetLocationSettingsCalls = append(ah.SetLocationSettingsCalls, input)
	return nil
}

func (ah testAppHandler) GetNetworkSettings() *models.NetworkSettings {
	return ah.NetworkSettings
}

func (ah testAppHandler) GetLocationSettings() *models.LocationSettings {
	return ah.LocationSettings
}

func (ah testAppHandler) GetAPSettings() *models.APSettings {
	return ah.APSettings
}

func (ah testAppHandler) GetTrackStatus() *models.TrackStatus {
	return ah.TrackStatus
}

func (ah *testAppHandler) SetAPSettings(input map[string]interface{}) error {
	ah.SetAPSettingsCalls = append(ah.SetAPSettingsCalls, input)
	return nil
}

func (ah testAppHandler) GetTime() *time.Time {
	v := time.Now()
	return &v
}

func (ah testAppHandler) GetFlags() *models.Flags {
	return ah.Flags
}

func (ah testAppHandler) GetAlignStatus() *models.AlignStatus {
	return ah.AlignStatus
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
