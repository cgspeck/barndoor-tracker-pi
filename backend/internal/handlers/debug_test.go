package handlers

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/pyros2097/cupaloy"
)

func TestDebugHandler(t *testing.T) {

	handler := newTestAppHandler()
	handler.H = DebugHandler

	handler.NetworkSettings = &models.NetworkSettings{
		AccessPointMode: true,
		APSettings: &models.APSettings{
			Channel: 11,
			Key:     "",
			SSID:    "barndoor-tracker",
		},
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
