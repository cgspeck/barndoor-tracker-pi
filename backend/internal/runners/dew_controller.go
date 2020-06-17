package runners

import (
	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type DewControllerRunner struct{}

func NewDewControllerRunner() (*DewControllerRunner, error) {
	return &DewControllerRunner{}, nil
}

func (dcr *DewControllerRunner) SetPID(p float64, i float64, d float64) {

}

func (dcr *DewControllerRunner) SetEnabled(enabled bool) {

}

func (dcr *DewControllerRunner) SetTargetTemperature(targetTemperature int) {

}

func (dcr *DewControllerRunner) GetStatus() *models.DewControllerStatus {
	return &models.DewControllerStatus{}
}
