package lsm9ds1

func magAvailable(axis Axis, status byte) bool {
	/*
		example of bitshift by four to disregard overrun flags
			dataAvail := status << 4

		example of disregard overrun flags
			dataAvail := status & 0x0F
	*/
	bAxis := byte(axis)
	axisVal := byte(0x01 << bAxis)
	axisTest := status & axisVal
	allTest := axisTest & 0x08
	return allTest == 0x08 || axisTest == axisVal
}

func gyroAvailable(status byte) bool {
	return status&0x02 == 0x02
}

func accelAvailable(status byte) bool {
	return status&0x01 == 0x01
}
