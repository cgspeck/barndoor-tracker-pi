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
