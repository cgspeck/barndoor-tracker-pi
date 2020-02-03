package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"

	"github.com/pyros2097/cupaloy"
)

func TestLocationSettingsHandler(t *testing.T) {
	handler := newTestAppHandler()
	handler.H = LocationSettingsHandler

	handler.LocationSettings = &models.LocationSettings{
		AltError:       1,
		AzError:        2,
		Latitude:       3.4,
		MagDeclination: -5.6,
		XOffset:        7,
		YOffset:        8,
		ZOffset:        9,
	}

	rr := doRequest(&handler, 200, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}

func doLocationSettingsPost(
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
func TestLocationSettingsHandlerPost(t *testing.T) {
	body := `{
	"Latitude": -37.4
}
`

	handler := newTestAppHandler()
	handler.H = LocationSettingsHandler
	handler.LocationSettings = &models.LocationSettings{
		AltError:          1,
		AzError:           2,
		Latitude:          3.4,
		MagDeclination:    -5.6,
		XOffset:           7,
		YOffset:           8,
		ZOffset:           9,
		ManagementEnabled: true,
	}

	rr := doLocationSettingsPost(body, &handler, http.StatusOK, t)

	// Check the response body is what we expect.
	err := cupaloy.Snapshot(rr)
	if err != nil {
		t.Error(err)
	}

	// if len(handler.SetAPModeCalls) != 1 {
	// 	t.Errorf("Expected call to SetAPMode")
	// }

	// if handler.SetAPModeCalls[0] != true {
	// 	t.Errorf("Expected true call to SetAPMode")
	// }
}

func TestLocationSettingsHandlerPostManagementDisabled(t *testing.T) {
	body := `{
	"Latitude": -37.4
}
`

	handler := newTestAppHandler()
	handler.H = LocationSettingsHandler
	handler.LocationSettings = &models.LocationSettings{
		AltError:          1,
		AzError:           2,
		Latitude:          3.4,
		MagDeclination:    -5.6,
		XOffset:           7,
		YOffset:           8,
		ZOffset:           9,
		ManagementEnabled: false,
	}

	rr := doLocationSettingsPost(body, &handler, http.StatusBadRequest, t)

	// Check the response body is what we expect.
	err := cupaloy.Snapshot(rr)
	if err != nil {
		t.Error(err)
	}

	// if len(handler.SetAPModeCalls) != 0 {
	// 	t.Errorf("Expected no call to SetAPMode")
	// }
}
