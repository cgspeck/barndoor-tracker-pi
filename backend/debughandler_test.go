package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pyros2097/cupaloy"
)

func TestDebugHandler(t *testing.T) {
	fmt.Println("hello testing!")

	req, err := http.NewRequest("GET", "/debug", nil)
	if err != nil {
		t.Fatal(err)
	}
	appContext := CreateAppContext(time.Time{})
	rr := httptest.NewRecorder()
	handler := AppHandler{appContext, DebugHandler}

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
