package pin_wrapper

import (
	_ "github.com/kidoman/embd/host/rpi"
)

type WrappedPin struct {
	pinNo int
}

func NewWrappedPin(pinNo int) (*WrappedPin, error) {

	return &WrappedPin{
		pinNo: pinNo,
	}, nil
}

func (wp *WrappedPin) SetHigh() {}

func (wp *WrappedPin) SetLow() {}
