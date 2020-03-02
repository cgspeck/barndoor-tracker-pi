package lsm9ds1

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
