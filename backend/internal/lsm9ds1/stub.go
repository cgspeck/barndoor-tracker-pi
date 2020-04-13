package lsm9ds1

/*
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
*/

type StubLSM9DS1 struct {
	AccelVal []float32
}

func (*StubLSM9DS1) Calibrate(bool)                {}
func (*StubLSM9DS1) CalibrateMag()                 {}
func (*StubLSM9DS1) LoadMagBias() error            { return nil }
func (*StubLSM9DS1) MagRange() []int16             { return []int16{} }
func (*StubLSM9DS1) GyroAvailable() bool           { return false }
func (*StubLSM9DS1) AccelAvailable() bool          { return false }
func (*StubLSM9DS1) MagAvailable(axis Axis) bool   { return false }
func (*StubLSM9DS1) ReadGyro() error               { return nil }
func (*StubLSM9DS1) ReadAccel() error              { return nil }
func (*StubLSM9DS1) ReadMag() error                { return nil }
func (*StubLSM9DS1) CalcGyro(gyro int16) float32   { return 0 }
func (*StubLSM9DS1) CalcAccel(accel int16) float32 { return 0 }
func (*StubLSM9DS1) CalcMag(mag int16) float32     { return 0 }

func (s *StubLSM9DS1) GetAccel() ([]float32, error) {
	return s.AccelVal, nil
}
