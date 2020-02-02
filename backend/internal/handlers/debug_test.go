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

func getWirelessInterfaceRpi() (string, error) {
	return "wlan0", nil
}

type fooError struct{}

func (e *fooError) Error() string {
	return "exit status 1"
}

func getWirelessInterfaceNone() (string, error) {
	return "", &fooError{}
}

func doDebugRequest(
	appHandlerFunc func(*models.AppContext, http.ResponseWriter, *http.Request) (int, error),
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

func TestDebugHandler(t *testing.T) {
	rr := doDebugRequest(DebugHandler, getWirelessInterfaceRpi, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}

func TestDebugHandlerNoWireless(t *testing.T) {
	rr := doDebugRequest(DebugHandler, getWirelessInterfaceNone, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}
