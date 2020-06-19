package runners

import (
	"fmt"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/pidlogger"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/cgspeck/barndoor-tracker-pi/internal/pin_wrapper"
	"github.com/felixge/pidctrl"
)

const samplePeriodSeconds = 10

type DewControllerRunner struct {
	pid                *pidctrl.PIDController
	enabled            bool
	processValue       float64
	previouslyEnabled  bool
	lastSampleTime     time.Time
	heatRequested      bool
	heatEntireDuration bool
	isHeatingNow       bool
	turnHeatOffAfter   time.Time
	heatPin            pin_wrapper.IWrappedPin
	pidLogger          *pidlogger.PIDLogger
	doLogging          bool
}

func NewDewControllerRunner(p float64, i float64, d float64, setPoint float64, enabled bool, heatPinNo int, doLogging bool, pidLogger *pidlogger.PIDLogger) (*DewControllerRunner, error) {
	pidctrl_ := pidctrl.NewPIDController(p, i, d)
	pidctrl_.SetOutputLimits(0, samplePeriodSeconds)
	pidctrl_.Set(setPoint)
	heatPin, err := pin_wrapper.NewWrappedPin(heatPinNo, false)

	if err != nil {
		return nil, err
	}

	return &DewControllerRunner{
		pid:               pidctrl_,
		enabled:           enabled,
		previouslyEnabled: enabled,
		processValue:      0,
		heatPin:           heatPin,
		doLogging:         doLogging,
		pidLogger:         pidLogger,
	}, nil
}

func (dcr *DewControllerRunner) SetPID(p float64, i float64, d float64) {
	dcr.pid.SetPID(p, i, d)
}

func (dcr *DewControllerRunner) SetEnabled(enabled bool) {
	dcr.enabled = enabled
}

func (dcr *DewControllerRunner) SetLoggingEnabled(enabled bool) {
	dcr.doLogging = enabled
}

func (dcr *DewControllerRunner) SetTargetTemperature(setPoint float64) {
	dcr.pid.Set(setPoint)
}

func (dcr *DewControllerRunner) GetStatus() *models.DewControllerStatus {
	p, i, d := dcr.pid.PID()
	return &models.DewControllerStatus{
		CurrentTemperature: dcr.processValue,
		CurrentlyHeating:   dcr.isHeatingNow,
		Enabled:            dcr.enabled,
		TargetTemperature:  dcr.pid.Get(),
		P:                  p,
		I:                  i,
		D:                  d,
		LastSampleTime:     dcr.lastSampleTime,
		HeatEntireDuration: dcr.heatEntireDuration,
		TurnHeatOffAfter:   dcr.turnHeatOffAfter,
		LoggingEnabled:     dcr.doLogging,
	}
}

func (dcr *DewControllerRunner) getProcessValue() float64 {
	return 0.0
}

func (dcr *DewControllerRunner) turnOffHeat() {
	dcr.heatPin.SetLow()
	dcr.isHeatingNow = false
}

func (dcr *DewControllerRunner) turnOnHeat() {
	dcr.heatPin.SetHigh()
	dcr.isHeatingNow = true
}

func (dcr *DewControllerRunner) Run(currentTime time.Time) {
	enabled := dcr.enabled
	previouslyEnabled := dcr.previouslyEnabled

	if previouslyEnabled && !enabled {
		dcr.turnOffHeat()
		dcr.heatRequested = false
		dcr.heatEntireDuration = false
	}

	if enabled {
		secondsSincePreviousCheck := currentTime.Sub(dcr.lastSampleTime).Seconds()

		if secondsSincePreviousCheck >= samplePeriodSeconds {
			pv := dcr.getProcessValue()
			hv := dcr.pid.Update(pv)
			dcr.lastSampleTime = currentTime

			if hv >= 0 {
				if int(hv) >= samplePeriodSeconds {
					dcr.heatEntireDuration = true
				} else {
					dcr.heatEntireDuration = false
					dcr.turnHeatOffAfter = currentTime.Add(time.Second * time.Duration(hv))
				}
				dcr.heatRequested = true
				dcr.turnOnHeat()
			} else {
				dcr.heatRequested = false
				dcr.heatEntireDuration = false
				dcr.turnOffHeat()
			}

			if dcr.doLogging {
				dcr.pidLogger.Emit(fmt.Sprintf("%v, %v, %v, %v\n", currentTime.Format(time.RFC3339), dcr.pid.Get(), pv, hv))
			}
		} else {
			if !dcr.heatEntireDuration && currentTime.After(dcr.turnHeatOffAfter) {
				dcr.turnOffHeat()
			}
		}

	}
	dcr.previouslyEnabled = enabled
}
