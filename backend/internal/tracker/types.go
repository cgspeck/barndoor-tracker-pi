package tracker

import "github.com/cgspeck/barndoor-tracker-pi/internal/mutexi2cbus"

type Tracker struct {
	status byte
	i2c    mutexi2cbus.I2CBus
}
