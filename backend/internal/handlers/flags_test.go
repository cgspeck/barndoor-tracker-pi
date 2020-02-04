package handlers

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/pyros2097/cupaloy"
)

func TestFlagHandler(t *testing.T) {
	handler := newTestAppHandler()
	handler.H = Flags

	handler.Flags = &models.Flags{
		NeedsLocationSettings: true,
		NeedsNetworkSettings:  false,
		RunningAsRoot:         true,
	}

	rr := doRequest(&handler, 200, t)

	err := cupaloy.Snapshot(rr)
	if err != nil {
		t.Error(err)
	}
}
