package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/models"

	"github.com/pyros2097/cupaloy"
)

func TestNetworkSettingsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/network_settings", nil)
	if err != nil {
		t.Fatal(err)
	}
	cmdFlags := models.CmdFlags{}

	appContext, err := config.CreateAppContext(time.Time{}, cmdFlags)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := AppHandler{appContext, NetworkSettingsHandler}

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}