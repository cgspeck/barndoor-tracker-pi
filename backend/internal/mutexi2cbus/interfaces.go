package mutexi2cbus

type I2CBus interface {
	// Read/Write a single byte from the given address
	ReadByteFromAddr(addr byte) (byte, error)
	WriteByteToAddr(addr, value byte) error
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
