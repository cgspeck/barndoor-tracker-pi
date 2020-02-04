package config

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/pyros2097/cupaloy"
)

func TestIsLocationConfigChangedEmptyInput(t *testing.T) {
	current := models.LocationSettings{
		AltError:          1,
		AzError:           2,
		Latitude:          3.4,
		MagDeclination:    -5.6,
		XOffset:           7,
		YOffset:           8,
		ZOffset:           9,
		ManagementEnabled: true,
	}
	input := make(map[string]interface{}, 0)

	changed, result, err := IsLocationConfigChanged(input, current)

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsLocationConfigChangedIrrelevantKeys(t *testing.T) {
	current := models.LocationSettings{
		AltError:          1,
		AzError:           2,
		Latitude:          3.4,
		MagDeclination:    -5.6,
		XOffset:           7,
		YOffset:           8,
		ZOffset:           9,
		ManagementEnabled: true,
	}
	input := make(map[string]interface{}, 0)
	input["foo"] = "bar"

	changed, result, err := IsLocationConfigChanged(input, current)

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsLocationConfigChangedSameValue(t *testing.T) {
	current := models.LocationSettings{
		AltError:          1,
		AzError:           2,
		Latitude:          3.4,
		MagDeclination:    -5.6,
		XOffset:           7,
		YOffset:           8,
		ZOffset:           9,
		ManagementEnabled: true,
	}
	input := make(map[string]interface{}, 0)
	input["latitude"] = 3.4
	input["zOffset"] = 9
	input["magDeclination"] = -5.6

	changed, result, err := IsLocationConfigChanged(input, current)

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsLocationConfigChangedValidDifferentValues(t *testing.T) {
	current := models.LocationSettings{
		AltError:          1,
		AzError:           2,
		Latitude:          3.4,
		MagDeclination:    -5.6,
		XOffset:           7,
		YOffset:           8,
		ZOffset:           9,
		ManagementEnabled: true,
	}
	input := make(map[string]interface{}, 0)
	input["latitude"] = -37.814
	input["zOffset"] = 9
	input["magDeclination"] = 11.72

	changed, result, err := IsLocationConfigChanged(input, current)

	err = cupaloy.Snapshot(changed, result, err)
	if err != nil {
		t.Error(err)
	}
}

func TestIsLocationConfigChangedInvalidDifferentValues(t *testing.T) {
	current := models.LocationSettings{
		AltError:          1,
		AzError:           2,
		Latitude:          3.4,
		MagDeclination:    -5.6,
		XOffset:           7,
		YOffset:           8,
		ZOffset:           9,
		ManagementEnabled: true,
	}
	input := make(map[string]interface{})
	input["latitude"] = -91.0

	changed, result, err := IsLocationConfigChanged(input, current)

	if !current.Equals(result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(LatitudeValidationError) {
		t.Error("Expected error")
	}

	input["latitude"] = 91.0
	changed, result, err = IsLocationConfigChanged(input, current)

	if !current.Equals(result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(LatitudeValidationError) {
		t.Error("Expected error")
	}

	input = make(map[string]interface{})
	input["magDeclination"] = 180.1
	changed, result, err = IsLocationConfigChanged(input, current)

	if !current.Equals(result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(MagneticDeclinationValidationError) {
		t.Error("Expected error")
	}

	input["magDeclination"] = -180.1
	changed, result, err = IsLocationConfigChanged(input, current)

	if !current.Equals(result) {
		t.Error("Did not expect settings to change")
	}
	if changed {
		t.Error("Did not expect changed flag to be true")
	}
	if err != err.(MagneticDeclinationValidationError) {
		t.Error("Expected error")
	}
}
