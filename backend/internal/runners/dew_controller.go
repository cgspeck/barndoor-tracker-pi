package runners

import (
	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/felixge/pidctrl"
)

const samplePeriodSeconds = 10

type DewControllerRunner struct {
	pid          *pidctrl.PIDController
	enabled      bool
	processValue float64
}

func NewDewControllerRunner(p float64, i float64, d float64, setPoint float64, enabled bool) (*DewControllerRunner, error) {
	pidctrl_ := pidctrl.NewPIDController(p, i, d)
	pidctrl_.SetOutputLimits(0, samplePeriodSeconds)
	pidctrl_.Set(setPoint)
	return &DewControllerRunner{
		pid:          pidctrl_,
		enabled:      enabled,
		processValue: 0,
	}, nil
}

func (dcr *DewControllerRunner) SetPID(p float64, i float64, d float64) {
	dcr.pid.SetPID(p, i, d)
}

func (dcr *DewControllerRunner) SetEnabled(enabled bool) {
	dcr.enabled = false
}

func (dcr *DewControllerRunner) SetTargetTemperature(setPoint float64) {
	dcr.pid.Set(setPoint)
}

func (dcr *DewControllerRunner) GetStatus() *models.DewControllerStatus {
	p, i, d := dcr.pid.PID()
	return &models.DewControllerStatus{
		CurrentTemperature: dcr.processValue,
		CurrentlyHeating:   false,
		Enabled:            dcr.enabled,
		P:                  p,
		I:                  i,
		D:                  d,
	}
}
