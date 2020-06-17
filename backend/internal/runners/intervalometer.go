package runners

import (
	"fmt"
	"sync"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/cgspeck/barndoor-tracker-pi/internal/pin_wrapper"
)

const idleCheckIntervalMillis = 1000
const focusDelayMillis = 1000

type IntervalometerMode int

const (
	IntervalometerEnabled = iota
	IntervalometerRunning
	IntervalometerDisabled
)

func (m IntervalometerMode) String() string {
	return [...]string{"Enabled", "Running", "Disabled"}[m]
}

type ShotPhase int

const (
	ShotPhaseRest = iota
	ShotPhaseFocussing
	ShotPhaseShooting
)

func (s ShotPhase) String() string {
	return [...]string{"Rest", "Focussing", "Shooting"}[s]
}

type IntervalometerRunner struct {
	sync.RWMutex
	focusPin        pin_wrapper.IWrappedPin
	shutterPin      pin_wrapper.IWrappedPin
	bulbTimeSeconds int
	restTimeSeconds int

	manageFocusPin bool
	mode           IntervalometerMode
	shootingPhase  ShotPhase

	nextAction time.Time
	shotCount  int
}

func NewIntervalometerRunner(shutterPinNo int, focusPinNo int, ip *models.IntervalPeriods) (*IntervalometerRunner, error) {
	shutterPin, err := pin_wrapper.NewWrappedPin(shutterPinNo)
	if err != nil {
		return nil, err
	}

	focusPin, err := pin_wrapper.NewWrappedPin(focusPinNo)
	if err != nil {
		return nil, err
	}

	return &IntervalometerRunner{
		shutterPin:      shutterPin,
		focusPin:        focusPin,
		bulbTimeSeconds: ip.BulbTimeSeconds,
		restTimeSeconds: ip.RestTimeSeconds,
	}, nil
}

func turnOffPin(pin int) {}

func turnOnPin(pin int) {}

func (ir *IntervalometerRunner) GetIntervalPeriods() models.IntervalPeriods {
	return models.IntervalPeriods{
		BulbTimeSeconds: ir.bulbTimeSeconds,
		RestTimeSeconds: ir.restTimeSeconds,
	}
}

func (ir *IntervalometerRunner) SetIntervalPeriods(ip *models.IntervalPeriods) {
	ir.bulbTimeSeconds = ip.BulbTimeSeconds
	ir.restTimeSeconds = ip.RestTimeSeconds
}

func (ir *IntervalometerRunner) doInitialFocusOrShoot(currentTime time.Time) (time.Time, ShotPhase, string) {
	var nextActionTime time.Time
	var nextShotPhase ShotPhase
	var tsNextIntervalometerStatus string

	if ir.manageFocusPin {
		ir.shotCount = 0
		ir.focusPin.SetHigh()
		nextActionTime = currentTime.Add(time.Duration(focusDelayMillis) * time.Millisecond)
		nextShotPhase = ShotPhaseFocussing
		tsNextIntervalometerStatus = nextShotPhase.String()
	} else {
		ir.shotCount = 1
		ir.shutterPin.SetHigh()
		nextActionTime = currentTime.Add(time.Duration(ir.bulbTimeSeconds) * time.Second)
		nextShotPhase = ShotPhaseShooting
		tsNextIntervalometerStatus = fmt.Sprintf("%s %d", nextShotPhase.String(), ir.shotCount)
	}

	return nextActionTime, nextShotPhase, tsNextIntervalometerStatus
}

func (ir *IntervalometerRunner) progressShoot(currentTime time.Time) (time.Time, ShotPhase, string) {
	var nextActionTime time.Time
	var nextShotPhase ShotPhase
	var tsNextIntervalometerStatus string
	switch ir.shootingPhase {
	case ShotPhaseRest:
		if ir.manageFocusPin {
			ir.focusPin.SetHigh()
			nextActionTime = currentTime.Add(time.Duration(focusDelayMillis) * time.Millisecond)
			nextShotPhase = ShotPhaseFocussing
			tsNextIntervalometerStatus = nextShotPhase.String()
		} else {
			ir.shotCount++
			ir.shutterPin.SetHigh()
			nextActionTime = currentTime.Add(time.Duration(ir.bulbTimeSeconds) * time.Second)
			nextShotPhase = ShotPhaseShooting
			tsNextIntervalometerStatus = fmt.Sprintf("%s %d", nextShotPhase.String(), ir.shotCount)
		}
	case ShotPhaseFocussing:
		ir.focusPin.SetLow()
		ir.shotCount++
		ir.shutterPin.SetHigh()
		nextActionTime = currentTime.Add(time.Duration(ir.bulbTimeSeconds) * time.Second)
		nextShotPhase = ShotPhaseShooting
		tsNextIntervalometerStatus = fmt.Sprintf("%s %d", nextShotPhase.String(), ir.shotCount)
	case ShotPhaseShooting:
		ir.shutterPin.SetLow()
		nextActionTime = currentTime.Add(time.Duration(ir.restTimeSeconds) * time.Second)
		nextShotPhase = ShotPhaseRest
		tsNextIntervalometerStatus = fmt.Sprintf("%s %d", nextShotPhase.String(), ir.shotCount)
	}
	return nextActionTime, nextShotPhase, tsNextIntervalometerStatus
}

func (ir *IntervalometerRunner) Run(currentTime time.Time, ts *models.TrackStatus) {
	if currentTime.After(ir.nextAction) {
		ts.Lock()
		defer ts.Unlock()
		ir.Lock()
		defer ir.Unlock()

		intervalometerShouldBeEnabled := ts.IntervalometerEnabled
		intervalometerShouldBeDisabled := !intervalometerShouldBeEnabled
		isTracking := ts.State == "Tracking"
		var nextActionTime time.Time
		var tsNextIntervalometerStatus string
		var nextShotPhase ShotPhase
		var nextIntervalometerMode IntervalometerMode

		if isTracking {
			switch ir.mode {
			case IntervalometerDisabled:
				if intervalometerShouldBeDisabled {
					break
				}
				nextIntervalometerMode = IntervalometerDisabled
				nextActionTime = currentTime.Add(time.Millisecond * idleCheckIntervalMillis)
				nextShotPhase = ShotPhaseRest
				tsNextIntervalometerStatus = nextIntervalometerMode.String()
			case IntervalometerEnabled:
				nextIntervalometerMode = IntervalometerRunning
				nextActionTime, nextShotPhase, tsNextIntervalometerStatus = ir.doInitialFocusOrShoot(currentTime)
			case IntervalometerRunning:
				nextIntervalometerMode = IntervalometerRunning
				nextActionTime, nextShotPhase, tsNextIntervalometerStatus = ir.progressShoot(currentTime)
			}
		} else {
			nextIntervalometerMode = ir.mode

			switch ir.mode {
			case IntervalometerDisabled:
				if intervalometerShouldBeDisabled {
					break
				}
				nextIntervalometerMode = IntervalometerDisabled
			case IntervalometerEnabled:
				if intervalometerShouldBeEnabled {
					break
				}
				nextIntervalometerMode = IntervalometerEnabled
			case IntervalometerRunning:
				ir.focusPin.SetLow()
				ir.shutterPin.SetLow()
				if intervalometerShouldBeEnabled {
					nextIntervalometerMode = IntervalometerEnabled
				} else {
					nextIntervalometerMode = IntervalometerDisabled
				}
			}
			nextActionTime = currentTime.Add(time.Millisecond * idleCheckIntervalMillis)
			tsNextIntervalometerStatus = nextIntervalometerMode.String()
		}

		ts.IntervalometerState = tsNextIntervalometerStatus
		ir.nextAction = nextActionTime
		ir.shootingPhase = nextShotPhase
		ir.mode = nextIntervalometerMode
	}
}
