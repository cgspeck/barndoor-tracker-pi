package runners

import (
	"fmt"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

const arduinoAddress = byte(0x04)
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

func (ir *IntervalometerRunner) Run(currentTime time.Time, ts *models.TrackStatus) {

	if currentTime.After(ir.nextAction) {
		ts.Lock()
		defer ts.Unlock()
		intervalometerShouldBeEnabled := ts.IntervolmeterEnabled
		intervalometerShouldRun := ts.State == "Tracking"
		var nextActionTime time.Time
		var nextIntervolmeterState int
		var tsNextIntervalometerStatus string
		var nextShotPhase ShotPhase
		var nextIntervalometerMode IntervalometerMode

		switch ir.mode {
		case IntervalometerRunning:
			if intervalometerShouldBeEnabled {
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
					turnOffPin(ir.focusPin)
					nextActionTime = currentTime.Add(time.Duration(ir.restTimeSeconds) * time.Second)
					nextShotPhase = ShotPhaseRest
					tsNextIntervalometerStatus = fmt.Sprintf("%s %d", nextShotPhase.String(), ir.shotCount)
				}
				nextIntervalometerMode = IntervalometerRunning
			} else {
				switch ir.shootingPhase {
				case ShotPhaseFocussing:
					turnOffPin(ir.focusPin)
				case ShotPhaseShooting:
					turnOffPin(ir.shutterPin)
				}
				ir.shotCount = 0
				nextShotPhase = ShotPhaseRest
				nextIntervalometerMode = IntervalometerDisabled
				nextActionTime = currentTime.Add(time.Millisecond * idleCheckIntervalMillis)
				tsNextIntervalometerStatus = ir.mode.String()
			}
		}
		ts.IntervolmeterState = tsNextIntervalometerStatus
		ir.nextAction = nextActionTime
		ir.shootingPhase = nextShotPhase
		ir.mode = nextIntervalometerMode
	}
}
