package mutexi2cbus

import (
	"sync"
)

type MutexI2cBus struct {
	mu sync.Mutex
}

func NewMutexI2cBus(bus int) MutexI2cBus {
	return MutexI2cBus{}
}

func (p *MutexI2cBus) ReadByteFromAddr(addr byte) (value byte, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return byte(0x00), nil
}

func (p *MutexI2cBus) WriteByteToAddr(addr, value byte) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return nil
}

// ReadFromReg reads n (len(value)) bytes from the given address and register.
func (p *MutexI2cBus) ReadFromReg(addr, reg byte, value []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return nil
}

// ReadByteFromReg reads a byte from the given address and register.
func (p *MutexI2cBus) ReadByteFromReg(addr, reg byte) (value byte, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	const WHO_AM_I_M = 0x0F
	const WHO_AM_I_XG = 0x0F
	const WHO_AM_I_AG_RSP = 0x68
	const WHO_AM_I_M_RSP = 0x3D

	if addr == byte(0x1e) && reg == WHO_AM_I_M {
		return WHO_AM_I_M_RSP, nil
	}

	if addr == byte(0x6b) && reg == WHO_AM_I_XG {
		return WHO_AM_I_AG_RSP, nil
	}

	return 0, nil
}

// WriteToReg writes len(value) bytes to the given address and register.
func (p *MutexI2cBus) WriteToReg(addr, reg byte, value []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return nil
}

// WriteByteToReg writes a byte to the given address and register.
func (p *MutexI2cBus) WriteByteToReg(addr, reg, value byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return nil
}

// WriteByteToReg writes a byte to the given address and register.
func (p *MutexI2cBus) WriteWordToReg(addr, reg byte, value uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return nil
}

// Close releases the resources associated with the bus.
func (p *MutexI2cBus) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return nil
}
