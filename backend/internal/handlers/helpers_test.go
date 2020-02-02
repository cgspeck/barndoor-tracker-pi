package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func getWirelessInterfaceRpi() (string, error) {
	return "wlan0", nil
}

func doRequest(
	appHandlerFunc func(IAppHandler, http.ResponseWriter, *http.Request) (int, error),
	getwirelessInterfaceFunc func() (string, error), t *testing.T) *httptest.ResponseRecorder {
	t.Helper()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cmdFlags := models.CmdFlags{}

	appContext, err := config.NewAppContext(time.Time{}, cmdFlags, getwirelessInterfaceFunc)
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
