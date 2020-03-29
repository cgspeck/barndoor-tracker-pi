package models

import (
	"reflect"
	"sync"
	"time"
)

// global config flags for frontend
type Flags struct {
	sync.RWMutex
	NeedsNetworkSettings  bool `json:"needsAPSettings"`
	NeedsLocationSettings bool `json:"needsLocationSettings"`
	RunningAsRoot         bool `json:"runningAsRoot"`
}

// configuration
type LocationSettings struct {
	sync.RWMutex
	Latitude          float64 `json:"latitude"`
	MagDeclination    float64 `json:"magDeclination"`
	AzError           float64 `json:"azError"`
	AltError          float64 `json:"altError"`
	XOffset           int     `json:"xOffset"`
	YOffset           int     `json:"yOffset"`
	ZOffset           int     `json:"zOffset"`
	ManagementEnabled bool    `json:"managementEnabled"`
}

func (l LocationSettings) Equals(o LocationSettings) bool {
	o.ManagementEnabled = l.ManagementEnabled
	return reflect.DeepEqual(l, o)
}

type APSettings struct {
	sync.RWMutex
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
	sync.RWMutex
	AccessPointMode   bool        `json:"accessPointMode"`
	APSettings        *APSettings `json:"aPSettings"`
	AvailableNetworks []*AvailableNetwork
	ManagementEnabled bool `json:"managementEnabled"`
	WirelessInterface string
	WirelessProfiles  []*WirelessProfile
}

// statuses
type AlignStatus struct {
	sync.RWMutex
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
	sync.RWMutex
	AlignStatus      *AlignStatus      `json:"alignStatus"`
	Flags            *Flags            `json:"flags"`
	LocationSettings *LocationSettings `json:"location"`
	Time             *time.Time        `json:"time"`
	NetworkSettings  *NetworkSettings  `json:"networkSettings"`
	OS               string            `json:"os"`
	Arch             string            `json:"arch"`
}
