package main

import "time"

// global config flags for frontend
type FlagStruct struct {
	NeedsNetworkSettings  bool
	NeedsLocationSettings bool
	IsRoot                bool
}

// configuration
type LocationStruct struct {
	Latitude       float64
	MagDeclination float64
	AzError        float64
	AltError       float64
	XOffset        int
	YOffset        int
	ZOffset        int
}

type APSettingsStruct struct {
	Channel int
	Key     string
	SSID    string
}

type WirelessStation struct {
	Key  string
	SSID string
}

type NetworkSettingsStruct struct {
	AccessPointMode   bool
	APSettings        *APSettingsStruct
	ManagementEnabled bool
	WirelessStations  []*WirelessStation
}

// statuses
type AlignStatusStruct struct {
	AzAligned  bool
	AltAligned bool
	CurrentAz  float64
	CurrentAlt float64
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
