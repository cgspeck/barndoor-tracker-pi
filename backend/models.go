package main

import "time"

// global config flags for frontend
type FlagStruct struct {
	NeedsNetworkSettings  bool
	NeedsLocationSettings bool
}

// configuration
type LocationStruct struct {
	Latitude       float32
	MagDeclination float32
	AzError        float32
	AltError       float32
	XOffset        int16
	YOffset        int16
	ZOffset        int16
}

type APSettingsStruct struct {
	SSID    string
	Key     string
	Channel int8
}

type WirelessStation struct {
	SSID string
	Key  string
}

type NetworkSettingsStruct struct {
	AccessPointMode  bool
	APSettings       *APSettingsStruct
	WirelessStations []*WirelessStation
}

// statuses
type AlignStatusStruct struct {
	AzAligned  bool
	AltAligned bool
	CurrentAz  float64
	CurrentAlt float32
}

// type TrackingStatus struct {
// 	Tracking            bool
// 	IntervolmeterEnable bool
// 	DewControlEnable    bool
// }

// the app context!
type AppContext struct {
	AlignStatus           *AlignStatusStruct
	Flags                 *FlagStruct
	Location              *LocationStruct
	PreviousTime          *time.Time
	NetworkSettingsStruct *NetworkSettingsStruct
}
