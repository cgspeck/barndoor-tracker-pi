package lsm9ds1

type MagnetoWhoamiFailed struct{}

func (_ MagnetoWhoamiFailed) Error() string {
	return "Magentometere Whoami Failed"
}

type AGWhoamiFailed struct{}

func (_ AGWhoamiFailed) Error() string {
	return "Accel/Gyro Whoami Failed"
}
