package models

import (
	"time"
)

// global config flags for frontend
type Flags struct {
	NeedsNetworkSettings  bool
	NeedsLocationSettings bool
	RunningAsRoot         bool
}

// configuration
type Location struct {
	Latitude       float64
	MagDeclination float64
	AzError        float64
	AltError       float64
	XOffset        int
	YOffset        int
	ZOffset        int
}

type APSettings struct {
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

type NetworkSettings struct {
	AccessPointMode bool
	*APSettings
	AvailableNetworks []*AvailableNetwork
	ManagementEnabled bool
	WirelessInterface string
	WirelessProfiles  []*WirelessProfile
}

// statuses
type AlignStatus struct {
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

// startup flags
type CmdFlags struct {
	DisableAP bool
}

// the app context!
type AppContext struct {
	AlignStatus     *AlignStatus     `json:"alignStatus"`
	Flags           *Flags           `json:"flags"`
	Location        *Location        `json:"location"`
	Time            *time.Time       `json:"time"`
	NetworkSettings *NetworkSettings `json:"networkSettings"`
}
