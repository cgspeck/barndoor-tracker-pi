package main

import (
	"time"
)

func CreateAppContext(previousTime time.Time) *AppContext {
	// TODO: load settings from configuration, if it exists
	var flags = FlagStruct{
		NeedsNetworkSettings:  true,
		NeedsLocationSettings: true,
	}

	var location = LocationStruct{
		Latitude:       -37.74,
		MagDeclination: 11.64,
		AzError:        2.0,
		AltError:       2.0,
		XOffset:        0,
		YOffset:        0,
		ZOffset:        0,
	}

	var alignStatus = AlignStatusStruct{
		AltAligned: true,
		AzAligned:  true,
		CurrentAz:  181.2,
		CurrentAlt: -37.4,
	}

	var networkSettings = NetworkSettingsStruct{
		AccessPointMode: true,
		APSettings: &APSettingsStruct{
			SSID:    "barndoor-tracker",
			Key:     "",
			Channel: 11,
		},
		WirelessStations: []*WirelessStation{},
	}
	return &AppContext{
		AlignStatus:           &alignStatus,
		Flags:                 &flags,
		Location:              &location,
		PreviousTime:          &previousTime,
		NetworkSettingsStruct: &networkSettings,
	}
}
