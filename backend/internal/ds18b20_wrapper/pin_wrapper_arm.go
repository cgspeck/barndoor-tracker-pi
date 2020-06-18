package ds18b20_wrapper

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

type WrappedPin struct {
	pinNo     int
	pin       embd.DigitalPin
	logicFlip bool
}

func NewWrappedPin(pinNo int, logicFlip bool) (*WrappedPin, error) {
	pin, err := embd.NewDigitalPin(pinNo)

	if err != nil {
		return nil, err
	}

	pin.SetDirection(embd.Out)
	pin.Write(embd.Low)

	return &WrappedPin{
		pinNo:     pinNo,
		pin:       pin,
		logicFlip: logicFlip,
	}, nil
}

func (wp *WrappedPin) SetHigh() {
	if wp.logicFlip {
		wp.pin.Write(embd.High)
	} else {
		wp.pin.Write(embd.Low)
	}

}

func (wp *WrappedPin) SetLow() {
	if wp.logicFlip {
		wp.pin.Write(embd.Low)
	} else {
		wp.pin.Write(embd.High)
	}
}
