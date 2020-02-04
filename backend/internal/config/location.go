package config

import (
	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type LatitudeValidationError struct{}

func (LatitudeValidationError) Error() string {
	return "invalid value"
}

type MagneticDeclinationValidationError struct{}

func (MagneticDeclinationValidationError) Error() string {
	return "invalid value"
}

// IsLocationConfigChanged takes a mapping of desired config and current config, and returns flag indiciating whether it changed and a copy of the new config
func IsLocationConfigChanged(
	input map[string]interface{},
	c models.LocationSettings,
) (bool, models.LocationSettings, error) {
	mustApplyChanges := false

	// 	"latitude", "magDeclination", "azError", "altError", "ZOffset", "yOffset", "zOffset",
	if lat, ok := input["latitude"].(float64); ok {

		if lat != c.Latitude {
			if lat > 90 || lat < -90 {
				return false, c, LatitudeValidationError{}
			}

			c.Latitude = lat
			mustApplyChanges = true
		}
	}

	if m, ok := input["magDeclination"].(float64); ok {
		if m != c.MagDeclination {
			if m > 180 || m < -180 {
				return false, c, MagneticDeclinationValidationError{}
			}
			c.MagDeclination = m
			mustApplyChanges = true
		}
	}

	if azError, ok := input["azError"].(float64); ok {
		if azError != c.AzError {
			c.AzError = azError
			mustApplyChanges = true
		}
	}

	if altError, ok := input["altError"].(float64); ok {
		if altError != c.AltError {
			c.AltError = altError
			mustApplyChanges = true
		}
	}

	if xOffset, ok := input["xOffset"].(int); ok {
		if xOffset != c.XOffset {
			c.XOffset = xOffset
			mustApplyChanges = true
		}
	}

	if yOffset, ok := input["yOffset"].(int); ok {
		if yOffset != c.YOffset {
			c.YOffset = yOffset
			mustApplyChanges = true
		}
	}

	if zOffset, ok := input["zOffset"].(int); ok {
		if zOffset != c.ZOffset {
			c.ZOffset = zOffset
			mustApplyChanges = true
		}
	}
	return mustApplyChanges, c, nil
}
