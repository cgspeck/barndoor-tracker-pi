package lsm9ds1

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

type ILSM9DS1 interface {
	New(I2CBus)
	Calibrate(bool)
	CalibrateMag()
	LoadMagBias() error
	MagRange() []int16
	GyroAvailable() bool
	AccelAvailable() bool
	MagAvailable(axis Axis) bool
	ReadGyro() error
	ReadAccel() error
	ReadMag() error
	CalcGyro(gyro int16) float32
	CalcAccel(accel int16) float32
	CalcMag(mag int16) float32
}
