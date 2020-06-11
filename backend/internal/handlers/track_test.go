package handlers

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"

	"github.com/pyros2097/cupaloy"
)

func TestTrackHandler(t *testing.T) {
	handler := newTestAppHandler()
	handler.H = TrackHandler

	handler.TrackStatus = &models.TrackStatus{
		State:              "foo",
		IntervolmeterState: "bar",
		ElapsedMillis:      123,
	}

	rr := doRequest(&handler, 200, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}

// func doAPSettingsPost(
// 	body string,
// 	testAppHandler *testAppHandler,
// 	expectedStatus int,
// 	t *testing.T) *httptest.ResponseRecorder {
// 	t.Helper()
// 	req, err := http.NewRequest("POST", "/", strings.NewReader(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()

// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	testAppHandler.ServeHTTP(rr, req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != expectedStatus {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, expectedStatus)
// 	}
// 	return rr
// }
// func TestAPSettingsHandlerPost(t *testing.T) {
// 	body := `{
// 	"key": "differentkey"
// }
// `

// 	handler := newTestAppHandler()
// 	handler.H = APSettingsHandler
// 	handler.APSettings = &models.APSettings{
// 		Channel: 11,
// 		Key:     "somekey1",
// 		SSID:    "somessid",
// 	}
// 	handler.NetworkSettings = &models.NetworkSettings{
// 		ManagementEnabled: true,
// 	}

// 	rr := doAPSettingsPost(body, &handler, http.StatusOK, t)

// 	if len(handler.SetAPSettingsCalls) != 1 {
// 		t.Errorf("Expected a call to SetAPSetting")
// 	}

// 	// Check the response body is what we expect.
// 	err := cupaloy.Snapshot(rr, handler.SetAPSettingsCalls[0])
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestAPSettingsHandlerPostManagementDisabled(t *testing.T) {
// 	body := `{
// 	"key": "differentkey"
// }
// `

// 	handler := newTestAppHandler()
// 	handler.H = APSettingsHandler
// 	handler.APSettings = &models.APSettings{
// 		Channel: 11,
// 		Key:     "somekey1",
// 		SSID:    "somessid",
// 	}
// 	handler.NetworkSettings = &models.NetworkSettings{
// 		ManagementEnabled: false,
// 	}

// 	rr := doAPSettingsPost(body, &handler, http.StatusBadRequest, t)

// 	// Check the response body is what we expect.
// 	err := cupaloy.Snapshot(rr)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(handler.SetAPSettingsCalls) != 0 {
// 		t.Errorf("Expected no calls to SetAPSetting")
// 	}
// }
