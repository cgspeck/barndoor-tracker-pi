package tracker

import (
	"github.com/cgspeck/barndoor-tracker-pi/internal/mutexi2cbus"
)

func New(i2c mutexi2cbus.I2CBus) (*Tracker, error) {
	l := Tracker{
		status: byte(0x00),
	}
	return &l, nil
}

func (t *Tracker) Poll() error {
	return nil
}

func (t *Tracker) Home() error {
	return nil
}

func (t *Tracker) Track() error {
	return nil
}

func (t *Tracker) Stop() error {
	return nil
}
