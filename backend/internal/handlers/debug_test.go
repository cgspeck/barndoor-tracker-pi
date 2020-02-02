package handlers

import (
	"testing"

	"github.com/pyros2097/cupaloy"
)

type fooError struct{}

func (e *fooError) Error() string {
	return "exit status 1"
}

func getWirelessInterfaceNone() (string, error) {
	return "", &fooError{}
}

func TestDebugHandler(t *testing.T) {
	rr := doRequest(DebugHandler, getWirelessInterfaceRpi, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}

func TestDebugHandlerNoWireless(t *testing.T) {
	rr := doRequest(DebugHandler, getWirelessInterfaceNone, t)

	// Check the response body is what we expect.
	err2 := cupaloy.Snapshot(rr)
	if err2 != nil {
		t.Error(err2)
	}
}
