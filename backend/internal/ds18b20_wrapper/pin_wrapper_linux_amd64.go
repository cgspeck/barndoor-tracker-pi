package ds18b20_wrapper

import (
	_ "github.com/kidoman/embd/host/rpi"
)

type WrappedPin struct {
	pinNo     int
	logicFlip bool
}

func NewWrappedPin(pinNo int, logicFlip bool) (*WrappedPin, error) {

	return &WrappedPin{
		pinNo:     pinNo,
		logicFlip: logicFlip,
	}, nil
}

func (wp *WrappedPin) SetHigh() {}

func (wp *WrappedPin) SetLow() {}
