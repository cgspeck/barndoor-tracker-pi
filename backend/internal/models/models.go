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
	Latitude       float64 `json:"latitude"`
	MagDeclination float64 `json:"magDeclination"`
	AzError        float64 `json:"azError"`
	AltError       float64 `json:"altError"`
	XOffset        int     `json:"xOffset"`
	YOffset        int     `json:"yOffset"`
	ZOffset        int     `json:"zOffset"`
	IgnoreAz       bool    `json:"ignoreAz"`
	IgnoreAlt      bool    `json:"ignoreAlt"`
}

func (l LocationSettings) Equals(o LocationSettings) bool {
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

type TrackStatus struct {
	sync.RWMutex
	State                 string    `json:"state"`
	PreviousState         string    `json:"previousState"`
	ElapsedMillis         int64     `json:"elapsedMillis"`
	IntervalometerEnabled bool      `json:"intervalometerEnabled"`
	IntervalometerState   string    `json:"intervalometerState"`
	TrackStartedAt        time.Time `json:"TrackStartedAt"`
}

type IntervalPeriods struct {
	sync.RWMutex
	BulbTimeSeconds int `json:"bulbInterval"`
	RestTimeSeconds int `json:"restInterval"`
}

type DewControllerSettings struct {
	sync.RWMutex
	TargetTemperature float64 `json:"targetTemperature"`
	Enabled           bool    `json:"dewControllerEnabled"`
	P                 float64 `json:"p"`
	I                 float64 `json:"i"`
	D                 float64 `json:"d"`
	LoggingEnabled    bool    `json:"loggingEnabled"`
}

type DewControllerStatus struct {
	TargetTemperature  float64   `json:"targetTemperature"`
	CurrentTemperature float64   `json:"currentTemperature"`
	CurrentlyHeating   bool      `json:"currentlyHeating"`
	Enabled            bool      `json:"dewControllerEnabled"`
	LoggingEnabled     bool      `json:"loggingEnabled"`
	P                  float64   `json:"p"`
	I                  float64   `json:"i"`
	D                  float64   `json:"d"`
	LastSampleTime     time.Time `json: "lastSampleTime"`
	HeatEntireDuration bool      `json: "heatEntireDuration"`
	TurnHeatOffAfter   time.Time `json: "turnOffHeatAfter"`
	DutyCycle          int       `json: "dutyCycle"`
	SensorOk           bool      `json:"sensorOk"`
}

type PIDLogFileRecord struct {
	Filename        string    `json:"filename"`
	EscapedFilename string    `json:"escapedFilename"`
	Size            int64     `json:"size"`
	Modified        time.Time `json:"modified"`
}

type PIDLogFiles struct {
	Files []PIDLogFileRecord `json:"files"`
}

func (ts *TrackStatus) ProcessTrackCommand(command string) (string, error) {
	nextState := ""
	currentState := ts.State
	stateChanged := false

	switch command {
	case "home":
		if currentState == "Idle" {
			stateChanged = true
			nextState = "Homing Requested"
		}
	case "track":
		if currentState == "Homed" {
			stateChanged = true
			nextState = "Tracking Requested"
		}
	case "stop":
		if currentState == "Tracking" {
			stateChanged = true
			nextState = "Stop Requested"
		}
	}

	if stateChanged {
		ts.PreviousState = currentState
		ts.State = nextState
		return nextState, nil
	}

	return "", InvalidStateChange{
		CurrentState: currentState,
		Requested:    command,
	}
}

func (ts *TrackStatus) ProcessArduinoStateChange(arduinoReportedState string) (string, error) {
	nextState := ""
	currentState := ts.State
	stateChanged := false

	switch arduinoReportedState {
	case "Homing":
		if currentState == "Idle" || currentState == "Homing Requested" {
			stateChanged = true
			nextState = "Homing"
		}
	case "Homed":
		if currentState == "Homing" || currentState == "Idle" {
			stateChanged = true
			nextState = "Homed"
		}
	case "Tracking":
		stateChanged = true
		nextState = "Tracking"
		ts.TrackStartedAt = time.Now()
	case "Idle":
		if currentState == "Tracking" || currentState == "Stop Requested" {
			stateChanged = true
			nextState = "Idle"
		}
	}

	if stateChanged {
		ts.PreviousState = currentState
		ts.State = nextState
		return nextState, nil
	}

	return "", InvalidStateChange{
		CurrentState: currentState,
		Requested:    arduinoReportedState,
	}
}

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
}

// the app context!
type AppContext struct {
	sync.RWMutex
	AlignStatus           *AlignStatus           `json:"alignStatus"`
	Flags                 *Flags                 `json:"flags"`
	LocationSettings      *LocationSettings      `json:"location"`
	Time                  *time.Time             `json:"time"`
	NetworkSettings       *NetworkSettings       `json:"networkSettings"`
	OS                    string                 `json:"os"`
	Arch                  string                 `json:"arch"`
	TrackStatus           *TrackStatus           `json:"trackStatus"`
	IntervalPeriods       *IntervalPeriods       `json:"intervalPeriods"`
	DewControllerSettings *DewControllerSettings `json: "dewControllerSettings"`
	PIDLogFiles           *PIDLogFiles           `json:"PIDLogFiles"`
}
