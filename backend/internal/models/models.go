package models

import (
	"time"
)

// global config flags for frontend
type Flags struct {
	NeedsNetworkSettings  bool `json:"needsAPSettings"`
	NeedsLocationSettings bool `json:"needsLocationSettings"`
	RunningAsRoot         bool
}

// configuration
type Location struct {
	Latitude          float64 `json:"latitude"`
	MagDeclination    float64 `json:"magDeclination"`
	AzError           float64 `json:"azError"`
	AltError          float64 `json:"altError"`
	XOffset           int     `json:"xOffset"`
	YOffset           int     `json:"yOffset"`
	ZOffset           int     `json:"zOffset"`
	ManagementEnabled bool    `json:"managementEnabled"`
}

type APSettings struct {
	Channel int    `json:"channel"`
	Key     string `json:"key"`
	SSID    string `json:"ssid"`
}

type WirelessProfile struct {
	Key  string `json:"key"`
	SSID string `json:"ssid"`
}

type AvailableNetwork struct {
	Channel     int
	Frequency   string
	SSID        string `json:"ssid"`
	SignalLevel int
}

type NetworkSettings struct {
	AccessPointMode   bool        `json:"accessPointMode"`
	APSettings        *APSettings `json:"aPSettings"`
	AvailableNetworks []*AvailableNetwork
	ManagementEnabled bool `json:"managementEnabled"`
	WirelessInterface string
	WirelessProfiles  []*WirelessProfile
}

// statuses
type AlignStatus struct {
	AzAligned  bool    `json:"azAligned"`
	AltAligned bool    `json:"altAligned"`
	CurrentAz  float64 `json:"currentAz"`
	CurrentAlt float64 `json:"currentAlt"`
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

type AppContextProvider interface {
	SetContext(a *AppContext)
	SetTime(currentTime time.Time)
	GetNetworkSettings() *NetworkSettings
	SetAPMode(bool)
	GetContext() *AppContext
	// WriteConfig() error
}

// the app context!
type AppContext struct {
	AlignStatus     *AlignStatus     `json:"alignStatus"`
	Flags           *Flags           `json:"flags"`
	Location        *Location        `json:"location"`
	Time            *time.Time       `json:"time"`
	NetworkSettings *NetworkSettings `json:"networkSettings"`
}
