package lsm9ds1

import (
	"sync"
)

type MutexI2cBus struct {
	mu sync.Mutex
}

func NewMutexI2cBus(bus int) MutexI2cBus {
	return MutexI2cBus{}
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
