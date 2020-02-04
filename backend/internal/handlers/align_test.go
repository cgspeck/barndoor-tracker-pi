package handlers

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/pyros2097/cupaloy"
)

func TestAlignHandler(t *testing.T) {

	handler := newTestAppHandler()
	handler.H = AlignHandler

	handler.AlignStatus = &models.AlignStatus{
		AltAligned: true,
		AzAligned:  true,
		CurrentAlt: -12.6,
		CurrentAz:  23.7,
	}

	rr := doRequest(&handler, 200, t)

	// Check the response body is what we expect.
	err := cupaloy.Snapshot(rr)
	if err != nil {
		t.Error(err)
	}
}
