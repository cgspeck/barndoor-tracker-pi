package config

import (
	"strings"
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/pyros2097/cupaloy"
)

func TestSaveConfig(t *testing.T) {
	var b strings.Builder

	c := configSettings{
		AccessPointMode: true,
		APSettings: &models.APSettings{
			Channel: 11,
			Key:     "some key",
			SSID:    "my amazing hotspot",
		},
		LocationSettings: &models.Location{
			AltError:       1,
			AzError:        2,
			Latitude:       3.4,
			MagDeclination: -5.6,
			XOffset:        7,
			YOffset:        8,
			ZOffset:        9,
		},
		NeedsNetworkSettings:  false,
		NeedsLocationSettings: true,
	}

	err := saveConfig(&c, &b)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = cupaloy.Snapshot(b.String())
	if err != nil {
		t.Error(err)
	}
}
