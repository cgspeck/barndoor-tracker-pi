package lsm9ds1

import (
	"sync"

	"github.com/kidoman/embd"
)

type I2CBus interface {
	// ReadFromReg reads n (len(value)) bytes from the given address and register.
	ReadFromReg(addr, reg byte, value []byte) error
	// ReadByteFromReg reads a byte from the given address and register.
	ReadByteFromReg(addr, reg byte) (value byte, err error)
	// WriteToReg writes len(value) bytes to the given address and register.
	WriteToReg(addr, reg byte, value []byte) error
	// WriteByteToReg writes a byte to the given address and register.
	WriteByteToReg(addr, reg, value byte) error
	// WriteU16ToReg
	WriteWordToReg(addr, reg byte, value uint16) error
	// Close releases the resources associated with the bus.
	Close() error
}

type MutexI2cBus struct {
	bus embd.I2CBus
	mu  sync.Mutex
}

func NewMutexI2cBus(bus int) MutexI2cBus {
	return MutexI2cBus{
		embd.NewI2CBus(byte(bus)),
		sync.Mutex{},
	}
}

// ReadFromReg reads n (len(value)) bytes from the given address and register.
func (p *MutexI2cBus) ReadFromReg(addr, reg byte, value []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.bus.ReadFromReg(addr, reg, value)
}

// ReadByteFromReg reads a byte from the given address and register.
func (p *MutexI2cBus) ReadByteFromReg(addr, reg byte) (value byte, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.bus.ReadByteFromReg(addr, reg)
}

// WriteToReg writes len(value) bytes to the given address and register.
func (p *MutexI2cBus) WriteToReg(addr, reg byte, value []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.bus.WriteToReg(addr, reg, value)
}

// WriteByteToReg writes a byte to the given address and register.
func (p *MutexI2cBus) WriteByteToReg(addr, reg, value byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.bus.WriteByteToReg(addr, reg, value)
}

// WriteByteToReg writes a byte to the given address and register.
func (p *MutexI2cBus) WriteWordToReg(addr, reg byte, value uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.bus.WriteWordToReg(addr, reg, value)
}

// Close releases the resources associated with the bus.
func (p *MutexI2cBus) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.bus.Close()
}
