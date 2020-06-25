package lsm9ds1

type ILSM9DS1 interface {
	Calibrate(bool)
	CalibrateMag()
	LoadMagBias() error
	MagRange() []int16
	GyroAvailable() bool
	AccelAvailable() bool
	MagAvailable(axis Axis) bool
	ReadGyro() error
	ReadAccel() error
	GetAccel() ([]float32, error)
	ReadMag() error
	CalcGyro(gyro int16) float32
	CalcAccel(accel int16) float32
	CalcMag(mag int16) float32
}
