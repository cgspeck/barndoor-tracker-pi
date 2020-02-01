package models

import "time"

// global config flags for frontend
type FlagStruct struct {
	NeedsNetworkSettings  bool
	NeedsLocationSettings bool
	RunningAsRoot         bool
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

type WirelessProfile struct {
	Key  string
	SSID string
}

type AvailableNetwork struct {
	Channel     int
	Frequency   string
	SSID        string
	SignalLevel int
}

type NetworkSettingsStruct struct {
	AccessPointMode   bool
	APSettings        *APSettingsStruct
	AvailableNetworks []*AvailableNetwork
	ManagementEnabled bool
	WirelessInterface string
	WirelessProfiles  []*WirelessProfile
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
	AlignStatus     *AlignStatusStruct
	Flags           *FlagStruct
	Location        *LocationStruct
	PreviousTime    *time.Time
	NetworkSettings *NetworkSettingsStruct
}
