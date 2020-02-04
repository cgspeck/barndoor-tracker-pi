package config

import (
	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type APChannelValidationError struct{}

func (APChannelValidationError) Error() string {
	return "invalid value"
}

type APKeyValidationError struct{}

func (APKeyValidationError) Error() string {
	return "invalid value"
}

type APSSIDValidationError struct{}

func (APSSIDValidationError) Error() string {
	return "invalid value"
}

// IsAPConfigChanged takes a mapping of desired config and current config, and returns flag indiciating whether it changed and a copy of the new config and an error
func IsAPConfigChanged(
	input map[string]interface{},
	c models.APSettings,
) (bool, models.APSettings, error) {
	mustApplyChanges := false

	// 	"channel" int, "key", "azError", "ssid"
	if channel, ok := input["channel"].(int); ok {
		if channel != c.Channel {

			if channel < 1 || channel > 14 {
				return false, c, APChannelValidationError{}
			}
			c.Channel = channel
			mustApplyChanges = true
		}
	}

	if key, ok := input["key"].(string); ok {
		if key != c.Key {
			if len(key) > 0 && len(key) < 8 {
				return false, c, APKeyValidationError{}
			}
			if len(key) > 63 {
				return false, c, APKeyValidationError{}
			}
			c.Key = key
			mustApplyChanges = true
		}
	}

	if ssid, ok := input["ssid"].(string); ok {
		if ssid != c.SSID {
			if len(ssid) > 32 {
				return false, c, APSSIDValidationError{}
			}
			c.SSID = ssid
			mustApplyChanges = true
		}
	}
	return mustApplyChanges, c, nil
}
