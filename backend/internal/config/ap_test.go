package config

import (
	"reflect"
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/pyros2097/cupaloy"
)

func TestIsAPConfigChangedEmptyInput(t *testing.T) {
	current := models.APSettings{
		Channel: 11,
		Key:     "",
		SSID:    "barndoor_tracker",
	}
	input := make(map[string]interface{}, 0)

	changed, result, err := IsAPConfigChanged(input, current)

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsAPConfigChangedIrrelevantKeys(t *testing.T) {
	current := models.APSettings{
		Channel: 11,
		Key:     "",
		SSID:    "barndoor_tracker",
	}
	input := make(map[string]interface{}, 0)
	input["foo"] = "bar"

	changed, result, err := IsAPConfigChanged(input, current)

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsAPConfigChangedSameValues(t *testing.T) {
	current := models.APSettings{
		Channel: 11,
		Key:     "",
		SSID:    "barndoor_tracker",
	}
	input := make(map[string]interface{}, 0)
	input["SSID"] = "barndoor_tracker"
	input["key"] = ""
	input["Channel"] = 11

	changed, result, err := IsAPConfigChanged(input, current)

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsAPConfigChangedValidDifferentValues(t *testing.T) {
	current := models.APSettings{
		Channel: 11,
		Key:     "",
		SSID:    "barndoor_tracker",
	}
	input := make(map[string]interface{}, 0)
	input["channel"] = 2
	input["key"] = "somekey1"
	input["ssid"] = "someSSID"

	changed, result, err := IsAPConfigChanged(input, current)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsAPConfigChangedSetBlankKey(t *testing.T) {
	current := models.APSettings{
		Channel: 11,
		Key:     "somekey",
		SSID:    "barndoor_tracker",
	}
	input := make(map[string]interface{}, 0)
	input["key"] = ""

	changed, result, err := IsAPConfigChanged(input, current)

	if result.Key != "" {
		t.Errorf("Expected blank key, got: %q", result.Key)
	}

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsAPConfigChangedInvalidDifferentValues(t *testing.T) {
	current := models.APSettings{
		Channel: 11,
		Key:     "",
		SSID:    "barndoor_tracker",
	}
	// Allowed channels: 1-14
	// Allowed SSID:
	// 	0-32 octets with arbitrary contents
	//   no character set associated with the SSID - a 32-byte string of NUL-bytes is a valid SSID
	// Key: passphrase must be a sequence of between 8 and 63 ASCII-encoded characters.
	input := make(map[string]interface{}, 0)
	input["channel"] = 0

	changed, result, err := IsAPConfigChanged(input, current)

	if !reflect.DeepEqual(current, result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(APChannelValidationError) {
		t.Error("Expected error")
	}

	input["channel"] = 15
	changed, result, err = IsAPConfigChanged(input, current)

	if !reflect.DeepEqual(current, result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(APChannelValidationError) {
		t.Error("Expected error")
	}

	input = make(map[string]interface{})
	input["key"] = " "
	changed, result, err = IsAPConfigChanged(input, current)

	if !reflect.DeepEqual(current, result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(APKeyValidationError) {
		t.Error("Expected error")
	}

	input["key"] = "12345678901234567890123456789012345678901234567890123456789012345"
	changed, result, err = IsAPConfigChanged(input, current)

	if !reflect.DeepEqual(current, result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(APKeyValidationError) {
		t.Error("Expected error")
	}

	input = make(map[string]interface{})
	input["ssid"] = "123456789012345678901234567890123"
	changed, result, err = IsAPConfigChanged(input, current)

	if !reflect.DeepEqual(current, result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(APSSIDValidationError) {
		t.Error("Expected error")
	}
}
