package runners

import (
	"fmt"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
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
	focusPin        int
	shutterPin      int
	bulbTimeSeconds int
	restTimeSeconds int

	manageFocusPin bool
	mode           IntervalometerMode
	shootingPhase  ShotPhase

	nextAction time.Time
	shotCount  int
}

func NewIntervalometerRunner(shutterPin int) *IntervalometerRunner {
	return &IntervalometerRunner{
		shutterPin: shutterPin,
	}
}

func turnOffPin(pin int) {}

func turnOnPin(pin int) {}

func (ir *IntervalometerRunner) setupRunState() {}

func (ir *IntervalometerRunner) setupEnabledState() {}

func (ir *IntervalometerRunner) setupDisabledState() {}

func (ir *IntervalometerRunner) doInitialFocusOrShoot(currentTime time.Time) (time.Time, ShotPhase, string) {
	var nextActionTime time.Time
	var nextShotPhase ShotPhase
	var tsNextIntervalometerStatus string

	if ir.manageFocusPin {
		ir.shotCount = 0
		turnOnPin(ir.focusPin)
		nextActionTime = currentTime.Add(time.Duration(focusDelayMillis) * time.Millisecond)
		nextShotPhase = ShotPhaseFocussing
		tsNextIntervalometerStatus = nextShotPhase.String()
	} else {
		ir.shotCount = 1
		turnOnPin(ir.shutterPin)
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
			turnOnPin(ir.focusPin)
			nextActionTime = currentTime.Add(time.Duration(focusDelayMillis) * time.Millisecond)
			nextShotPhase = ShotPhaseFocussing
			tsNextIntervalometerStatus = nextShotPhase.String()
		} else {
			ir.shotCount++
			turnOnPin(ir.shutterPin)
			nextActionTime = currentTime.Add(time.Duration(ir.bulbTimeSeconds) * time.Second)
			nextShotPhase = ShotPhaseShooting
			tsNextIntervalometerStatus = fmt.Sprintf("%s %d", nextShotPhase.String(), ir.shotCount)
		}
	case ShotPhaseFocussing:
		turnOffPin(ir.focusPin)
		ir.shotCount++
		turnOnPin(ir.shutterPin)
		nextActionTime = currentTime.Add(time.Duration(ir.bulbTimeSeconds) * time.Second)
		nextShotPhase = ShotPhaseShooting
		tsNextIntervalometerStatus = fmt.Sprintf("%s %d", nextShotPhase.String(), ir.shotCount)
	case ShotPhaseShooting:
		turnOffPin(ir.shutterPin)
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
		intervalometerShouldBeEnabled := ts.IntervolmeterEnabled
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
				turnOffPin(ir.focusPin)
				turnOffPin(ir.shutterPin)
				if intervalometerShouldBeEnabled {
					nextIntervalometerMode = IntervalometerEnabled
				} else {
					nextIntervalometerMode = IntervalometerDisabled
				}
			}
			nextActionTime = currentTime.Add(time.Millisecond * idleCheckIntervalMillis)
			tsNextIntervalometerStatus = nextIntervalometerMode.String()
		}

		ts.IntervolmeterState = tsNextIntervalometerStatus
		ir.nextAction = nextActionTime
		ir.shootingPhase = nextShotPhase
		ir.mode = nextIntervalometerMode
	}
}
