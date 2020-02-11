package lsm9ds1

import (
	"log"

	"github.com/kidoman/embd"
)

/******************************************************************************
SFE_LSM9DS1 Library Source File

Ported from library by [SparkFun Electronics](https://github.com/sparkfun/LSM9DS1_Breakout)

This file implements all functions of the LSM9DS1 class. Functions here range
from higher level stuff, like reading/writing LSM9DS1 registers to low-level,
hardware reads and writes. I2C handler functions can be found towards the bottom of this file.

******************************************************************************/

// #include <Wire.h> // Wire library is used for I2C
// #include <SPI.h>  // SPI library is used for...SPI.

// Sensor Sensitivity Constants
// Values set according to the typical specifications provided in
// table 3 of the LSM9DS1 datasheet. (pg 12)
const SENSITIVITY_ACCELEROMETER_2 = 0.000061
const SENSITIVITY_ACCELEROMETER_4 = 0.000122
const SENSITIVITY_ACCELEROMETER_8 = 0.000244
const SENSITIVITY_ACCELEROMETER_16 = 0.000732
const SENSITIVITY_GYROSCOPE_245 = 0.00875
const SENSITIVITY_GYROSCOPE_500 = 0.0175
const SENSITIVITY_GYROSCOPE_2000 = 0.07
const SENSITIVITY_MAGNETOMETER_4 = 0.00014
const SENSITIVITY_MAGNETOMETER_8 = 0.00029
const SENSITIVITY_MAGNETOMETER_12 = 0.00043
const SENSITIVITY_MAGNETOMETER_16 = 0.00058

func (l *LSM9DS1) init() {
	l.settings.gyro.enabled = true
	l.settings.gyro.enableX = true
	l.settings.gyro.enableY = true
	l.settings.gyro.enableZ = true
	// gyro scale can be 245, 500, or 2000
	l.settings.gyro.scale = 245
	// gyro sample rate: value between 1-6
	// 1 = 14.9    4 = 238
	// 2 = 59.5    5 = 476
	// 3 = 119     6 = 952
	l.settings.gyro.sampleRate = 6
	// gyro cutoff frequency: value between 0-3
	// Actual value of cutoff frequency depends
	// on sample rate.
	l.settings.gyro.bandwidth = 0
	l.settings.gyro.lowPowerEnable = false
	l.settings.gyro.HPFEnable = false
	// Gyro HPF cutoff frequency: value between 0-9
	// Actual value depends on sample rate. Only applies
	// if gyroHPFEnable is true.
	l.settings.gyro.HPFCutoff = 0
	l.settings.gyro.flipX = false
	l.settings.gyro.flipY = false
	l.settings.gyro.flipZ = false
	l.settings.gyro.orientation = 0
	l.settings.gyro.latchInterrupt = true

	l.settings.accel.enabled = true
	l.settings.accel.enableX = true
	l.settings.accel.enableY = true
	l.settings.accel.enableZ = true
	// accel scale can be 2, 4, 8, or 16
	l.settings.accel.scale = 2
	// accel sample rate can be 1-6
	// 1 = 10 Hz    4 = 238 Hz
	// 2 = 50 Hz    5 = 476 Hz
	// 3 = 119 Hz   6 = 952 Hz
	l.settings.accel.sampleRate = 6
	// Accel cutoff freqeuncy can be any value between -1 - 3.
	// -1 = bandwidth determined by sample rate
	// 0 = 408 Hz   2 = 105 Hz
	// 1 = 211 Hz   3 = 50 Hz
	l.settings.accel.bandwidth = -1
	l.settings.accel.highResEnable = false
	// accelHighResBandwidth can be any value between 0-3
	// LP cutoff is set to a factor of sample rate
	// 0 = ODR/50    2 = ODR/9
	// 1 = ODR/100   3 = ODR/400
	l.settings.accel.highResBandwidth = 0

	l.settings.mag.enabled = true
	// mag scale can be 4, 8, 12, or 16
	l.settings.mag.scale = 4
	// mag data rate can be 0-7
	// 0 = 0.625 Hz  4 = 10 Hz
	// 1 = 1.25 Hz   5 = 20 Hz
	// 2 = 2.5 Hz    6 = 40 Hz
	// 3 = 5 Hz      7 = 80 Hz
	l.settings.mag.sampleRate = 7
	l.settings.mag.tempCompensationEnable = false
	// magPerformance can be any value between 0-3
	// 0 = Low power mode      2 = high performance
	// 1 = medium performance  3 = ultra-high performance
	l.settings.mag.XYPerformance = 3
	l.settings.mag.ZPerformance = 3
	l.settings.mag.lowPowerEnable = false
	// magOperatingMode can be 0-2
	// 0 = continuous conversion
	// 1 = single-conversion
	// 2 = power down
	l.settings.mag.operatingMode = 0

	l.settings.temp.enabled = true
	for i := 0; i < 3; i++ {
		l.gBias[i] = 0
		l.aBias[i] = 0
		l.mBias[i] = 0
		l.gBiasRaw[i] = 0
		l.aBiasRaw[i] = 0
		l.mBiasRaw[i] = 0
	}
	l._autoCalc = false
}

func (l *LSM9DS1) begin(agAddress uint, mAddress uint, wirePort *embd.I2CBus) uint {
	// Set device settings, they are used in many other places
	l.settings.device.commInterface = IMU_MODE_I2C
	l.settings.device.agAddress = agAddress
	l.settings.device.mAddress = mAddress
	l.settings.device.i2c = wirePort

	//! Todo: don't use _xgAddress or _mAddress, duplicating memory
	_xgAddress := l.settings.device.agAddress
	_mAddress := l.settings.device.mAddress

	l.init()

	l.constrainScales()
	// Once we have the scale values, we can calculate the resolution
	// of each sensor. That's what these functions are for. One for each sensor
	l.calcgRes() // Calculate DPS / ADC tick, stored in gRes variable
	l.calcmRes() // Calculate Gs / ADC tick, stored in mRes variable
	l.calcaRes() // Calculate g / ADC tick, stored in aRes variable

	// We expect caller to begin their I2C port, with the speed of their choice external to the library
	// But if they forget, we could start the hardware here.
	// l.settings.device.i2c->begin();	// Initialize I2C library

	// To verify communication, we can read from the WHO_AM_I register of
	// each device. Store those in a variable so we can return them.
	mTest := l.mReadByte(WHO_AM_I_M)    // Read the gyro WHO_AM_I
	xgTest := l.xgReadByte(WHO_AM_I_XG) // Read the accel/mag WHO_AM_I
	whoAmICombined := (xgTest << 8) | mTest

	if whoAmICombined != ((WHO_AM_I_AG_RSP << 8) | WHO_AM_I_M_RSP) {
		return 0
	}
	// Gyro initialization stuff:
	l.initGyro() // This will "turn on" the gyro. Setting up interrupts, etc.

	// Accelerometer initialization stuff:
	l.initAccel() // "Turn on" all axes of the accel. Set up interrupts, etc.

	// Magnetometer initialization stuff:
	l.initMag() // "Turn on" all axes of the mag. Set up interrupts, etc.

	// Once everything is initialized, return the WHO_AM_I registers we read:
	return whoAmICombined
}

func (l *LSM9DS1) beginSPI(ag_CS_pin uint, m_CS_pin uint) uint {
	// Set device settings, they are used in many other places
	l.settings.device.commInterface = IMU_MODE_SPI
	l.settings.device.agAddress = ag_CS_pin
	l.settings.device.mAddress = m_CS_pin

	//! Todo: don't use _xgAddress or _mAddress, duplicating memory
	_xgAddress := l.settings.device.agAddress
	_mAddress := l.settings.device.mAddress

	l.init()

	l.constrainScales()
	// Once we have the scale values, we can calculate the resolution
	// of each sensor. That's what these functions are for. One for each sensor
	l.calcgRes() // Calculate DPS / ADC tick, stored in gRes variable
	l.calcmRes() // Calculate Gs / ADC tick, stored in mRes variable
	l.calcaRes() // Calculate g / ADC tick, stored in aRes variable

	// Now, initialize our hardware interface.
	l.initSPI() // Initialize SPI

	// To verify communication, we can read from the WHO_AM_I register of
	// each device. Store those in a variable so we can return them.
	mTest := l.mReadByte(WHO_AM_I_M)    // Read the gyro WHO_AM_I
	xgTest := l.xgReadByte(WHO_AM_I_XG) // Read the accel/mag WHO_AM_I
	whoAmICombined := (xgTest << 8) | mTest

	if whoAmICombined != ((WHO_AM_I_AG_RSP << 8) | WHO_AM_I_M_RSP) {
		return 0
	}

	// Gyro initialization stuff:
	l.initGyro() // This will "turn on" the gyro. Setting up interrupts, etc.

	// Accelerometer initialization stuff:
	l.initAccel() // "Turn on" all axes of the accel. Set up interrupts, etc.

	// Magnetometer initialization stuff:
	l.initMag() // "Turn on" all axes of the mag. Set up interrupts, etc.

	// Once everything is initialized, return the WHO_AM_I registers we read:
	return whoAmICombined
}

func (l *LSM9DS1) initGyro() {
	var tempRegValue uint = 0

	// CTRL_REG1_G (Default value: 0x00)
	// [ODR_G2][ODR_G1][ODR_G0][FS_G1][FS_G0][0][BW_G1][BW_G0]
	// ODR_G[2:0] - Output data rate selection
	// FS_G[1:0] - Gyroscope full-scale selection
	// BW_G[1:0] - Gyroscope bandwidth selection

	// To disable gyro, set sample rate bits to 0. We'll only set sample
	// rate if the gyro is enabled.
	if l.settings.gyro.enabled {
		tempRegValue = uint((l.settings.gyro.sampleRate & 0x07) << 5)
	}

	switch l.settings.gyro.scale {
	case 500:
		tempRegValue |= (0x1 << 3)
	case 2000:
		tempRegValue |= (0x3 << 3)
		// Otherwise we'll set it to 245 dps (0x0 << 4)
	}

	tempRegValue |= (l.settings.gyro.bandwidth & 0x3)
	l.xgWriteByte(CTRL_REG1_G, tempRegValue)

	// CTRL_REG2_G (Default value: 0x00)
	// [0][0][0][0][INT_SEL1][INT_SEL0][OUT_SEL1][OUT_SEL0]
	// INT_SEL[1:0] - INT selection configuration
	// OUT_SEL[1:0] - Out selection configuration
	l.xgWriteByte(CTRL_REG2_G, 0x00)

	// CTRL_REG3_G (Default value: 0x00)
	// [LP_mode][HP_EN][0][0][HPCF3_G][HPCF2_G][HPCF1_G][HPCF0_G]
	// LP_mode - Low-power mode enable (0: disabled, 1: enabled)
	// HP_EN - HPF enable (0:disabled, 1: enabled)
	// HPCF_G[3:0] - HPF cutoff frequency
	if l.settings.gyro.lowPowerEnable {
		tempRegValue = 1 << 7
	} else {
		tempRegValue = 0
	}

	if l.settings.gyro.HPFEnable {
		tempRegValue |= (1 << 6) | (l.settings.gyro.HPFCutoff & 0x0F)
	}
	l.xgWriteByte(CTRL_REG3_G, tempRegValue)

	// CTRL_REG4 (Default value: 0x38)
	// [0][0][Zen_G][Yen_G][Xen_G][0][LIR_XL1][4D_XL1]
	// Zen_G - Z-axis output enable (0:disable, 1:enable)
	// Yen_G - Y-axis output enable (0:disable, 1:enable)
	// Xen_G - X-axis output enable (0:disable, 1:enable)
	// LIR_XL1 - Latched interrupt (0:not latched, 1:latched)
	// 4D_XL1 - 4D option on interrupt (0:6D used, 1:4D used)
	tempRegValue = 0
	if l.settings.gyro.enableZ {
		tempRegValue |= (1 << 5)
	}
	if l.settings.gyro.enableY {
		tempRegValue |= (1 << 4)
	}
	if l.settings.gyro.enableX {
		tempRegValue |= (1 << 3)
	}
	if l.settings.gyro.latchInterrupt {
		tempRegValue |= (1 << 1)
	}
	l.xgWriteByte(CTRL_REG4, tempRegValue)

	// ORIENT_CFG_G (Default value: 0x00)
	// [0][0][SignX_G][SignY_G][SignZ_G][Orient_2][Orient_1][Orient_0]
	// SignX_G - Pitch axis (X) angular rate sign (0: positive, 1: negative)
	// Orient [2:0] - Directional user orientation selection
	tempRegValue = 0
	if l.settings.gyro.flipX {
		tempRegValue |= (1 << 5)
	}
	if l.settings.gyro.flipY {
		tempRegValue |= (1 << 4)
	}
	if l.settings.gyro.flipZ {
		tempRegValue |= (1 << 3)
	}
	l.xgWriteByte(ORIENT_CFG_G, tempRegValue)
}

func (l *LSM9DS1) initAccel() {
	var tempRegValue uint = 0

	//	CTRL_REG5_XL (0x1F) (Default value: 0x38)
	//	[DEC_1][DEC_0][Zen_XL][Yen_XL][Zen_XL][0][0][0]
	//	DEC[0:1] - Decimation of accel data on OUT REG and FIFO.
	//		00: None, 01: 2 samples, 10: 4 samples 11: 8 samples
	//	Zen_XL - Z-axis output enabled
	//	Yen_XL - Y-axis output enabled
	//	Xen_XL - X-axis output enabled
	if l.settings.accel.enableZ {
		tempRegValue |= (1 << 5)
	}
	if l.settings.accel.enableY {
		tempRegValue |= (1 << 4)
	}
	if l.settings.accel.enableX {
		tempRegValue |= (1 << 3)
	}

	l.xgWriteByte(CTRL_REG5_XL, tempRegValue)

	// CTRL_REG6_XL (0x20) (Default value: 0x00)
	// [ODR_XL2][ODR_XL1][ODR_XL0][FS1_XL][FS0_XL][BW_SCAL_ODR][BW_XL1][BW_XL0]
	// ODR_XL[2:0] - Output data rate & power mode selection
	// FS_XL[1:0] - Full-scale selection
	// BW_SCAL_ODR - Bandwidth selection
	// BW_XL[1:0] - Anti-aliasing filter bandwidth selection
	tempRegValue = 0
	// To disable the accel, set the sampleRate bits to 0.
	if l.settings.accel.enabled {
		tempRegValue |= (l.settings.accel.sampleRate & 0x07) << 5
	}
	switch l.settings.accel.scale {
	case 4:
		tempRegValue |= (0x2 << 3)
	case 8:
		tempRegValue |= (0x3 << 3)
	case 16:
		tempRegValue |= (0x1 << 3)
		// Otherwise it'll be set to 2g (0x0 << 3)
	}

	if l.settings.accel.bandwidth >= 0 {
		tempRegValue |= (1 << 2) // Set BW_SCAL_ODR
		tempRegValue |= uint(l.settings.accel.bandwidth & 0x03)
	}
	l.xgWriteByte(CTRL_REG6_XL, tempRegValue)

	// CTRL_REG7_XL (0x21) (Default value: 0x00)
	// [HR][DCF1][DCF0][0][0][FDS][0][HPIS1]
	// HR - High resolution mode (0: disable, 1: enable)
	// DCF[1:0] - Digital filter cutoff frequency
	// FDS - Filtered data selection
	// HPIS1 - HPF enabled for interrupt function
	tempRegValue = 0
	if l.settings.accel.highResEnable {
		tempRegValue |= (1 << 7) // Set HR bit
		tempRegValue |= (l.settings.accel.highResBandwidth & 0x3) << 5
	}
	l.xgWriteByte(CTRL_REG7_XL, tempRegValue)
}

// This is a function that uses the FIFO to accumulate sample of accelerometer and gyro data, average
// them, scales them to  gs and deg/s, respectively, and then passes the biases to the main sketch
// for subtraction from all subsequent data. There are no gyro and accelerometer bias registers to store
// the data as there are in the ADXL345, a precursor to the LSM9DS0, or the MPU-9150, so we have to
// subtract the biases ourselves. This results in a more accurate measurement in general and can
// remove errors due to imprecise or varying initial placement. Calibration of sensor data in this manner
// is good practice.
func (l *LSM9DS1) calibrate(autoCalc bool) {
	var samples uint = 0
	var ii uint

	aBiasRawTemp := [3]int{0, 0, 0}
	gBiasRawTemp := [3]int{0, 0, 0}

	// Turn on FIFO and set threshold to 32 samples
	l.enableFIFO(true)
	l.setFIFO(FIFO_THS, 0x1F)
	for samples < 0x1F {
		samples = (l.xgReadByte(FIFO_SRC) & 0x3F) // Read number of stored samples
	}
	for ii = 0; ii < samples; ii++ {
		// Read the gyro data stored in the FIFO
		l.readGyro()
		gBiasRawTemp[0] += l.gx
		gBiasRawTemp[1] += l.gy
		gBiasRawTemp[2] += l.gz
		l.readAccel()
		aBiasRawTemp[0] += l.ax
		aBiasRawTemp[1] += l.ay
		aBiasRawTemp[2] += l.az - int(1./l.aRes) // Assumes sensor facing up!
	}

	for ii = 0; ii < 3; ii++ {
		l.gBiasRaw[ii] = gBiasRawTemp[ii] / int(samples)
		l.gBias[ii] = l.calcGyro(l.gBiasRaw[ii])
		l.aBiasRaw[ii] = aBiasRawTemp[ii] / int(samples)
		l.aBias[ii] = l.calcAccel(l.aBiasRaw[ii])
	}

	l.enableFIFO(false)
	l.setFIFO(FIFO_OFF, 0x00)

	if autoCalc {
		l._autoCalc = true
	}
}

func (l *LSM9DS1) calibrateMag(loadIn bool) {
	var i, j int
	magMin := [3]int{0, 0, 0}
	magMax := [3]int{0, 0, 0} // The road warrior

	for i = 0; i < 128; i++ {
		for l.magAvailable(X_AXIS) == 0 {
		}
		l.readMag()
		magTemp := [3]int{0, 0, 0}
		magTemp[0] = l.mx
		magTemp[1] = l.my
		magTemp[2] = l.mz
		for j = 0; j < 3; j++ {
			if magTemp[j] > magMax[j] {
				magMax[j] = magTemp[j]
			}
			if magTemp[j] < magMin[j] {
				magMin[j] = magTemp[j]
			}
		}
	}
	for j = 0; j < 3; j++ {
		l.mBiasRaw[j] = (magMax[j] + magMin[j]) / 2
		l.mBias[j] = l.calcMag(l.mBiasRaw[j])
		if loadIn {
			l.magOffset(Axis(j), l.mBiasRaw[j])
		}
	}

}

func (l *LSM9DS1) magOffset(axis Axis, offset int) {
	if axis > 2 {
		return
	}
	var msb, lsb int
	msb = (offset & 0xFF00) >> 8
	lsb = offset & 0x00FF
	iAxis := int(axis)
	l.miWriteByte(OFFSET_X_REG_L_M+(2*iAxis), lsb)
	l.miWriteByte(OFFSET_X_REG_H_M+(2*iAxis), msb)
}

func (l *LSM9DS1) initMag() {
	var tempRegValue uint = 0

	// CTRL_REG1_M (Default value: 0x10)
	// [TEMP_COMP][OM1][OM0][DO2][DO1][DO0][0][ST]
	// TEMP_COMP - Temperature compensation
	// OM[1:0] - X & Y axes op mode selection
	//	00:low-power, 01:medium performance
	//	10: high performance, 11:ultra-high performance
	// DO[2:0] - Output data rate selection
	// ST - Self-test enable
	if l.settings.mag.tempCompensationEnable {
		tempRegValue |= (1 << 7)
	}
	tempRegValue |= (l.settings.mag.XYPerformance & 0x3) << 5
	tempRegValue |= (l.settings.mag.sampleRate & 0x7) << 2
	l.mWriteByte(CTRL_REG1_M, tempRegValue)

	// CTRL_REG2_M (Default value 0x00)
	// [0][FS1][FS0][0][REBOOT][SOFT_RST][0][0]
	// FS[1:0] - Full-scale configuration
	// REBOOT - Reboot memory content (0:normal, 1:reboot)
	// SOFT_RST - Reset config and user registers (0:default, 1:reset)
	tempRegValue = 0
	switch l.settings.mag.scale {
	case 8:
		tempRegValue |= (0x1 << 5)
		break
	case 12:
		tempRegValue |= (0x2 << 5)
		break
	case 16:
		tempRegValue |= (0x3 << 5)
		break
		// Otherwise we'll default to 4 gauss (00)
	}
	l.mWriteByte(CTRL_REG2_M, tempRegValue) // +/-4Gauss

	// CTRL_REG3_M (Default value: 0x03)
	// [I2C_DISABLE][0][LP][0][0][SIM][MD1][MD0]
	// I2C_DISABLE - Disable I2C interace (0:enable, 1:disable)
	// LP - Low-power mode cofiguration (1:enable)
	// SIM - SPI mode selection (0:write-only, 1:read/write enable)
	// MD[1:0] - Operating mode
	//	00:continuous conversion, 01:single-conversion,
	//  10,11: Power-down
	tempRegValue = 0
	if l.settings.mag.lowPowerEnable {
		tempRegValue |= (1 << 5)
	}
	tempRegValue |= (l.settings.mag.operatingMode & 0x3)
	l.mWriteByte(CTRL_REG3_M, tempRegValue) // Continuous conversion mode

	// CTRL_REG4_M (Default value: 0x00)
	// [0][0][0][0][OMZ1][OMZ0][BLE][0]
	// OMZ[1:0] - Z-axis operative mode selection
	//	00:low-power mode, 01:medium performance
	//	10:high performance, 10:ultra-high performance
	// BLE - Big/little endian data
	tempRegValue = 0
	tempRegValue = (l.settings.mag.ZPerformance & 0x3) << 2
	l.mWriteByte(CTRL_REG4_M, tempRegValue)

	// CTRL_REG5_M (Default value: 0x00)
	// [0][BDU][0][0][0][0][0][0]
	// BDU - Block data update for magnetic data
	//	0:continuous, 1:not updated until MSB/LSB are read
	tempRegValue = 0
	l.mWriteByte(CTRL_REG5_M, tempRegValue)
}

func (l *LSM9DS1) accelAvailable() uint {
	status := l.xgReadByte(STATUS_REG_1)

	return (status & (1 << 0))
}

func (l *LSM9DS1) gyroAvailable() uint {
	status := l.xgReadByte(STATUS_REG_1)

	return ((status & (1 << 1)) >> 1)
}

func (l *LSM9DS1) tempAvailable() uint {
	status := l.xgReadByte(STATUS_REG_1)

	return ((status & (1 << 2)) >> 2)
}

func (l *LSM9DS1) magAvailable(axis Axis) uint {
	status := l.mReadByte(STATUS_REG_M)

	return ((status & (1 << axis)) >> axis)
}

func (l *LSM9DS1) readAccel() {
	// We'll read six bytes from the accelerometer into temp
	temp := [6]int{}

	// Read 6 bytes, beginning at OUT_X_L_XL
	if l.xgReadBytes(OUT_X_L_XL, &temp, 6) == 6 {
		l.ax = (temp[1] << 8) | temp[0] // Store x-axis values into ax
		l.ay = (temp[3] << 8) | temp[2] // Store y-axis values into ay
		l.az = (temp[5] << 8) | temp[4] // Store z-axis values into az
		if l._autoCalc {
			l.ax -= l.aBiasRaw[X_AXIS]
			l.ay -= l.aBiasRaw[Y_AXIS]
			l.az -= l.aBiasRaw[Z_AXIS]
		}
	}
}

func (l *LSM9DS1) readAccelAxis(axis axis) int {
	uaxis := uint(axis)
	temp := make([]int, 0)
	var value int
	if l.xgReadBytes(OUT_X_L_XL+(2*uaxis), &temp, 2) == 2 {
		value = (temp[1] << 8) | temp[0]

		if _autoCalc {
			value -= l.aBiasRaw[axis]
		}

		return value
	}
	return 0
}

func (l *LSM9DS1) readMag() {
	temp := [6]uint{} // We'll read six bytes from the mag into temp

	// Read 6 bytes, beginning at OUT_X_L_M
	if mReadBytes(OUT_X_L_M, temp, 6) == 6 {
		l.mx = (temp[1] << 8) | temp[0] // Store x-axis values into mx
		l.my = (temp[3] << 8) | temp[2] // Store y-axis values into my
		l.mz = (temp[5] << 8) | temp[4] // Store z-axis values into mz
	}
}

func (l *LSM9DS1) readMagAxis(axis uint) int {
	temp := [2]uint{}
	if mReadBytes(OUT_X_L_M+(2*axis), temp, 2) == 2 {
		return (temp[1] << 8) | temp[0]
	}
	return 0
}

func (l *LSM9DS1) readTemp() {
	temp := [2]int{} // We'll read two bytes from the temperature sensor into temp
	// Read 2 bytes, beginning at OUT_TEMP_L
	if l.xgReadBytes(OUT_TEMP_L, temp, 2) == 2 {
		offset := 25 // Per datasheet sensor outputs 0 typically @ 25 degrees centigrade
		temperature = offset + (((temp[1] << 8) | temp[0]) >> 8)
	}
}

func (l *LSM9DS1) readGyro() {
	temp := [6]uint{} // We'll read six bytes from the gyro into temp
	// Read 6 bytes, beginning at OUT_X_L_G
	if l.xgReadBytes(OUT_X_L_G, temp, 6) == 6 {
		l.gx = (temp[1] << 8) | temp[0] // Store x-axis values into gx
		l.gy = (temp[3] << 8) | temp[2] // Store y-axis values into gy
		l.gz = (temp[5] << 8) | temp[4] // Store z-axis values into gz
		if _autoCalc {
			l.gx -= l.gBiasRaw[X_AXIS]
			l.gy -= l.gBiasRaw[Y_AXIS]
			l.gz -= l.gBiasRaw[Z_AXIS]
		}
	}
}

func (l *LSM9DS1) readGyroAxis(axis uint) int {
	temp := [2]int{}
	var value int

	if l.xgReadBytes(OUT_X_L_G+(2*axis), temp, 2) == 2 {
		value = (temp[1] << 8) | temp[0]

		if _autoCalc {
			value -= l.gBiasRaw[axis]
		}
		return value
	}
	return 0
}

func (l *LSM9DS1) calcGyro(gyro int) float32 {
	// Return the gyro raw reading times our pre-calculated DPS / (ADC tick):
	return l.gRes * gyro
}

func (l *LSM9DS1) calcAccel(accel int) float32 {
	// Return the accel raw reading times our pre-calculated g's / (ADC tick):
	return l.aRes * accel
}

func (l *LSM9DS1) calcMag(mag int) float32 {
	// Return the mag raw reading times our pre-calculated Gs / (ADC tick):
	return l.mRes * mag
}

func (l *LSM9DS1) setGyroScale(gScl uint) {
	// Read current value of CTRL_REG1_G:
	ctrl1RegValue := l.xgReadByte(CTRL_REG1_G)
	// Mask out scale bits (3 & 4):
	ctrl1RegValue &= 0xE7
	switch gScl; {
	case 500:
		ctrl1RegValue |= (0x1 << 3)
		l.settings.gyro.scale = 500
	case 2000:
		ctrl1RegValue |= (0x3 << 3)
		l.settings.gyro.scale = 2000
	default: // Otherwise we'll set it to 245 dps (0x0 << 4)
		l.settings.gyro.scale = 245
	}
	l.xgWriteByte(CTRL_REG1_G, ctrl1RegValue)

	calcgRes()
}

func (l *LSM9DS1) setAccelScale(aScl uint) {
	// We need to preserve the other bytes in CTRL_REG6_XL. So, first read it:
	tempRegValue := l.xgReadByte(CTRL_REG6_XL)
	// Mask out accel scale bits:
	tempRegValue &= 0xE7

	switch aScl {
	case 4:
		tempRegValue |= (0x2 << 3)
		l.settings.accel.scale = 4
	case 8:
		tempRegValue |= (0x3 << 3)
		l.settings.accel.scale = 8
	case 16:
		tempRegValue |= (0x1 << 3)
		l.settings.accel.scale = 16
	default: // Otherwise it'll be set to 2g (0x0 << 3)
		l.settings.accel.scale = 2
	}
	l.xgWriteByte(CTRL_REG6_XL, tempRegValue)

	// Then calculate a new aRes, which relies on aScale being set correctly:
	calcaRes()
}

func (l *LSM9DS1) setMagScale(mScl uint) {
	// We need to preserve the other bytes in CTRL_REG6_XM. So, first read it:
	temp := mReadByte(CTRL_REG2_M)
	// Then mask out the mag scale bits:
	temp &= 0xFF ^ (0x3 << 5)

	switch mScl {
	case 8:
		temp |= (0x1 << 5)
		l.settings.mag.scale = 8
		break
	case 12:
		temp |= (0x2 << 5)
		l.settings.mag.scale = 12
		break
	case 16:
		temp |= (0x3 << 5)
		l.settings.mag.scale = 16
		break
	default: // Otherwise we'll default to 4 gauss (00)
		l.settings.mag.scale = 4
		break
	}

	// And write the new register value back into CTRL_REG6_XM:
	mWriteByte(CTRL_REG2_M, temp)

	// We've updated the sensor, but we also need to update our class variables
	// First update mScale:
	//mScale = mScl
	// Then calculate a new mRes, which relies on mScale being set correctly:
	calcmRes()
}

func (l *LSM9DS1) setGyroODR(gRate uint) {
	// Only do this if gRate is not 0 (which would disable the gyro)
	if (gRate & 0x07) != 0 {
		// We need to preserve the other bytes in CTRL_REG1_G. So, first read it:
		temp := l.xgReadByte(CTRL_REG1_G)
		// Then mask out the gyro ODR bits:
		temp &= 0xFF ^ (0x7 << 5)
		temp |= (gRate & 0x07) << 5
		// Update our settings struct
		l.settings.gyro.sampleRate = gRate & 0x07
		// And write the new register value back into CTRL_REG1_G:
		l.xgWriteByte(CTRL_REG1_G, temp)
	}
}

func (l *LSM9DS1) setAccelODR(aRate uint) {
	// Only do this if aRate is not 0 (which would disable the accel)
	if (aRate & 0x07) != 0 {
		// We need to preserve the other bytes in CTRL_REG1_XM. So, first read it:
		temp := l.xgReadByte(CTRL_REG6_XL)
		// Then mask out the accel ODR bits:
		temp &= 0x1F
		// Then shift in our new ODR bits:
		temp |= ((aRate & 0x07) << 5)
		l.settings.accel.sampleRate = aRate & 0x07
		// And write the new register value back into CTRL_REG1_XM:
		l.xgWriteByte(CTRL_REG6_XL, temp)
	}
}

func (l *LSM9DS1) setMagODR(mRate uint) {
	// We need to preserve the other bytes in CTRL_REG5_XM. So, first read it:
	temp := mReadByte(CTRL_REG1_M)
	// Then mask out the mag ODR bits:
	temp &= 0xFF ^ (0x7 << 2)
	// Then shift in our new ODR bits:
	temp |= ((mRate & 0x07) << 2)
	l.settings.mag.sampleRate = mRate & 0x07
	// And write the new register value back into CTRL_REG5_XM:
	mWriteByte(CTRL_REG1_M, temp)
}

func (l *LSM9DS1) calcgRes() {
	switch l.settings.gyro.scale; {
	case 245:
		l.gRes = SENSITIVITY_GYROSCOPE_245

	case 500:
		l.gRes = SENSITIVITY_GYROSCOPE_500

	case 2000:
		l.gRes = SENSITIVITY_GYROSCOPE_2000
	}
}

func (l *LSM9DS1) calcaRes() {
	switch l.settings.accel.scale; {
	case 2:
		l.aRes = SENSITIVITY_ACCELEROMETER_2

	case 4:
		l.aRes = SENSITIVITY_ACCELEROMETER_4

	case 8:
		l.aRes = SENSITIVITY_ACCELEROMETER_8

	case 16:
		l.aRes = SENSITIVITY_ACCELEROMETER_16

	}
}

func (l *LSM9DS1) calcmRes() {
	switch l.settings.mag.scale; {
	case 4:
		l.mRes = SENSITIVITY_MAGNETOMETER_4

	case 8:
		l.mRes = SENSITIVITY_MAGNETOMETER_8

	case 12:
		l.mRes = SENSITIVITY_MAGNETOMETER_12

	case 16:
		l.mRes = SENSITIVITY_MAGNETOMETER_16

	}
}

func (l *LSM9DS1) configInt(interrupt uint, generator uint,
	activeLow bool, pushPull bool) {
	// Write to INT1_CTRL or INT2_CTRL. [interupt] should already be one of
	// those two values.
	// [generator] should be an OR'd list of values from the interrupt_generators enum
	l.xgWriteByte(interrupt, generator)

	// Configure CTRL_REG8
	var temp uint
	temp = l.xgReadByte(CTRL_REG8)

	if activeLow {
		temp |= (1 << 5)
	} else {
		temp &= ^(1 << 5)
	}

	if pushPull {
		temp &= ^(1 << 4)
	} else {
		temp |= (1 << 4)
	}

	l.xgWriteByte(CTRL_REG8, temp)
}

func (l *LSM9DS1) configInactivity(duration uint, threshold uint, sleepOn bool) {
	temp := 0

	temp = threshold & 0x7F
	if sleepOn {
		temp |= (1 << 7)
	}
	l.xgWriteByte(ACT_THS, temp)

	l.xgWriteByte(ACT_DUR, duration)
}

func (l *LSM9DS1) getInactivity() uint {
	temp := l.xgReadByte(STATUS_REG_0)
	temp &= (0x10)
	return temp
}

func (l *LSM9DS1) configAccelInt(generator uint, andInterrupts bool) {
	// Use variables from accel_interrupt_generator, OR'd together to create
	// the [generator]value.
	temp := generator
	if andInterrupts {
		temp |= 0x80
	}
	l.xgWriteByte(INT_GEN_CFG_XL, temp)
}

func (l *LSM9DS1) configAccelThs(threshold uint, axis uint, duration uint, wait bool) {
	// Write threshold value to INT_GEN_THS_?_XL.
	// axis will be 0, 1, or 2 (x, y, z respectively)
	l.xgWriteByte(INT_GEN_THS_X_XL+axis, threshold)

	// Write duration and wait to INT_GEN_DUR_XL
	var temp uint
	temp = (duration & 0x7F)
	if wait {
		temp |= 0x80
	}
	l.xgWriteByte(INT_GEN_DUR_XL, temp)
}

func (l *LSM9DS1) getAccelIntSrc() uint {
	intSrc := l.xgReadByte(INT_GEN_SRC_XL)

	// Check if the IA_XL (interrupt active) bit is set
	if intSrc & (1 << 6) {
		return (intSrc & 0x3F)
	}

	return 0
}

func (l *LSM9DS1) configGyroInt(generator uint, aoi bool, latch bool) {
	// Use variables from accel_interrupt_generator, OR'd together to create
	// the [generator]value.
	temp := generator
	if aoi {
		temp |= 0x80
	}
	if latch {
		temp |= 0x40
	}
	l.xgWriteByte(INT_GEN_CFG_G, temp)
}

func (l *LSM9DS1) configGyroThs(threshold int, axis uint, duration uint, wait bool) {
	buffer := [2]uint{}
	buffer[0] = (threshold & 0x7F00) >> 8
	buffer[1] = (threshold & 0x00FF)
	// Write threshold value to INT_GEN_THS_?H_G and  INT_GEN_THS_?L_G.
	// axis will be 0, 1, or 2 (x, y, z respectively)
	l.xgWriteByte(INT_GEN_THS_XH_G+(axis*2), buffer[0])
	l.xgWriteByte(INT_GEN_THS_XH_G+1+(axis*2), buffer[1])

	// Write duration and wait to INT_GEN_DUR_XL
	var temp uint
	temp = (duration & 0x7F)
	if wait {
		temp |= 0x80
	}
	l.xgWriteByte(INT_GEN_DUR_G, temp)
}

func (l *LSM9DS1) getGyroIntSrc() uint {
	intSrc := l.xgReadByte(INT_GEN_SRC_G)

	// Check if the IA_G (interrupt active) bit is set
	if intSrc & (1 << 6) {
		return (intSrc & 0x3F)
	}

	return 0
}

func (l *LSM9DS1) configMagInt(generator uint, activeLow uint, latch bool) {
	// Mask out non-generator bits (0-4)
	config := (generator & 0xE0)
	// IEA bit is 0 for active-low, 1 for active-high.
	if activeLow == INT_ACTIVE_HIGH {
		config |= (1 << 2)
	}
	// IEL bit is 0 for latched, 1 for not-latched
	if !latch {
		config |= (1 << 1)
	}
	// As long as we have at least 1 generator, enable the interrupt
	if generator != 0 {
		config |= (1 << 0)
	}

	mWriteByte(INT_CFG_M, config)
}

func (l *LSM9DS1) configMagThs(threshold uint) {
	// Write high eight bits of [threshold] to INT_THS_H_M
	mWriteByte(INT_THS_H_M, uint((threshold&0x7F00)>>8))
	// Write low eight bits of [threshold] to INT_THS_L_M
	mWriteByte(INT_THS_L_M, uint(threshold&0x00FF))
}

func (l *LSM9DS1) getMagIntSrc() uint {
	intSrc := mReadByte(INT_SRC_M)

	// Check if the INT (interrupt active) bit is set
	if intSrc & (1 << 0) {
		return (intSrc & 0xFE)
	}

	return 0
}

func (l *LSM9DS1) sleepGyro(enable bool) {
	temp := l.xgReadByte(CTRL_REG9)
	if enable {
		temp |= (1 << 6)
	} else {
		temp &= ^(1 << 6)
	}
	l.xgWriteByte(CTRL_REG9, temp)
}

func (l *LSM9DS1) enableFIFO(enable bool) {
	temp := l.xgReadByte(CTRL_REG9)
	if enable {
		temp |= (1 << 1)
	} else {
		temp &= ^(1 << 1)
	}
	l.xgWriteByte(CTRL_REG9, temp)
}

func (l *LSM9DS1) setFIFO(fifoMode uint, fifoThs uint) {
	// Limit threshold - 0x1F (31) is the maximum. If more than that was asked
	// limit it to the maximum.
	var threshold uint
	if fifoThs <= 0x1F {
		threshold = fifoThs
	} else {
		threshold = 0x1F
	}
	l.xgWriteByte(FIFO_CTRL, ((fifoMode&0x7)<<5)|(threshold&0x1F))
}

func (l *LSM9DS1) getFIFOSamples() uint {
	return (l.xgReadByte(FIFO_SRC) & 0x3F)
}

func (l *LSM9DS1) constrainScales() {
	if (l.settings.gyro.scale != 245) && (l.settings.gyro.scale != 500) && (l.settings.gyro.scale != 2000) {
		l.settings.gyro.scale = 245
	}

	if (l.settings.accel.scale != 2) && (l.settings.accel.scale != 4) && (l.settings.accel.scale != 8) && (l.settings.accel.scale != 16) {
		l.settings.accel.scale = 2
	}

	if (l.settings.mag.scale != 4) && (l.settings.mag.scale != 8) && (l.settings.mag.scale != 12) && (l.settings.mag.scale != 16) {
		l.settings.mag.scale = 4
	}
}

func (l *LSM9DS1) xgWriteByte(subAddress uint, data uint) {
	// Whether we're using I2C or SPI, write a byte using the
	// gyro-specific I2C address or SPI CS pin.
	if l.settings.device.commInterface == IMU_MODE_I2C {
		I2CwriteByte(_xgAddress, subAddress, data)
	} else if l.settings.device.commInterface == IMU_MODE_SPI {
		// SPIwriteByte(_xgAddress, subAddress, data)
	}
}

func (l *LSM9DS1) miWriteByte(subAddress int, data int) {
	return mWriteByte(uint(subAddress), uint(data))
}

func (l *LSM9DS1) mWriteByte(subAddress uint, data uint) {
	// Whether we're using I2C or SPI, write a byte using the
	// accelerometer-specific I2C address or SPI CS pin.
	if l.settings.device.commInterface == IMU_MODE_I2C {
		return I2CwriteByte(_mAddress, subAddress, data)
	} else if l.settings.device.commInterface == IMU_MODE_SPI {
		// return SPIwriteByte(_mAddress, subAddress, data)
	}
}

func (l *LSM9DS1) xgReadByte(subAddress uint) uint {
	// Whether we're using I2C or SPI, read a byte using the
	// gyro-specific I2C address or SPI CS pin.
	if l.settings.device.commInterface == IMU_MODE_I2C {
		return I2CreadByte(_xgAddress, subAddress)
	} else if l.settings.device.commInterface == IMU_MODE_SPI {
		// return SPIreadByte(_xgAddress, subAddress)
	}
	return -1
}

func (l *LSM9DS1) xgReadBytes(subAddress uint, dest *[]uint, count uint) uint {
	// Whether we're using I2C or SPI, read multiple bytes using the
	// gyro-specific I2C address or SPI CS pin.
	if l.settings.device.commInterface == IMU_MODE_I2C {
		return I2CreadBytes(_xgAddress, subAddress, dest, count)
	} else if l.settings.device.commInterface == IMU_MODE_SPI {
		// return SPIreadBytes(_xgAddress, subAddress, dest, count)
	}
	return -1
}

func (l *LSM9DS1) mReadByte(subAddress uint) uint {
	uint
	// Whether we're using I2C or SPI, read a byte using the
	// accelerometer-specific I2C address or SPI CS pin.
	if l.settings.device.commInterface == IMU_MODE_I2C {
		return I2CreadByte(_mAddress, subAddress)
	} else if l.settings.device.commInterface == IMU_MODE_SPI {
		// return SPIreadByte(_mAddress, subAddress)
	}
	return -1
}

func (l *LSM9DS1) mReadBytes(subAddress uint, dest *uint, count uint) uint {
	// Whether we're using I2C or SPI, read multiple bytes using the
	// accelerometer-specific I2C address or SPI CS pin.
	if l.settings.device.commInterface == IMU_MODE_I2C {
		return I2CreadBytes(_mAddress, subAddress, dest, count)
	} else if l.settings.device.commInterface == IMU_MODE_SPI {
		// return SPIreadBytes(_mAddress, subAddress, dest, count)
	}
	return -1
}

func (l *LSM9DS1) initSPI() {
	pinMode(_xgAddress, OUTPUT)
	digitalWrite(_xgAddress, HIGH)
	pinMode(_mAddress, OUTPUT)
	digitalWrite(_mAddress, HIGH)

	SPI.begin()
	// Maximum SPI frequency is 10MHz, could divide by 2 here:
	SPI.setClockDivider(SPI_CLOCK_DIV2)
	// Data is read and written MSb first.
	SPI.setBitOrder(MSBFIRST)
	// Data is captured on rising edge of clock (CPHA = 0)
	// Base value of the clock is HIGH (CPOL = 1)
	SPI.setDataMode(SPI_MODE0)
}

/*
func (l *LSM9DS1) SPIwriteByte(csPin uint, subAddress uint, data uint) {
	digitalWrite(csPin, LOW) // Initiate communication

	// If write, bit 0 (MSB) should be 0
	// If single write, bit 1 should be 0
	SPI.transfer(subAddress & 0x3F) // Send Address
	SPI.transfer(data)              // Send data

	digitalWrite(csPin, HIGH) // Close communication
}

func (l *LSM9DS1) SPIreadByte(csPin uint, subAddress uint) uint {
	var temp uint
	// Use the multiple read function to read 1 byte.
	// Value is returned to `temp`.
	SPIreadBytes(csPin, subAddress, &temp, 1)
	return temp
}

func (l *LSM9DS1) SPIreadBytes(
	csPin uint,
	subAddress uint,
	dest *uint,
	count uint) {
	// To indicate a read, set bit 0 (msb) of first byte to 1
	rAddress := 0x80 | (subAddress & 0x3F)
	// Mag SPI port is different. If we're reading multiple bytes,
	// set bit 1 to 1. The remaining six bytes are the address to be read
	if (csPin == _mAddress) && count > 1 {
		rAddress |= 0x40
	}
	digitalWrite(csPin, LOW) // Initiate communication
	SPI.transfer(rAddress)

	for i := 0; i < count; i++ {
		dest[i] = SPI.transfer(0x00) // Read into destination array
	}
	digitalWrite(csPin, HIGH) // Close communication

	return count
}
*/
// Wire.h read and write protocols
func (l *LSM9DS1) I2CwriteByte(address int, subAddress int, data int) {
	l.settings.device.i2c.beginTransmission(address) // Initialize the Tx buffer
	l.settings.device.i2c.write(subAddress)          // Put slave register address in Tx buffer
	l.settings.device.i2c.write(data)                // Put data in Tx buffer
	l.settings.device.i2c.endTransmission()          // Send the Tx buffer
}

func (l *LSM9DS1) I2CreadByte(address uint, subAddress uint) uint {
	var data uint // `data` will store the register data

	l.settings.device.i2c.beginTransmission(address)     // Initialize the Tx buffer
	l.settings.device.i2c.write(subAddress)              // Put slave register address in Tx buffer
	l.settings.device.i2c.endTransmission(false)         // Send the Tx buffer, but send a restart to keep connection alive
	l.settings.device.i2c.requestFrom(address, 1.(uint)) // Read one byte from slave register address

	data = l.settings.device.i2c.read() // Fill Rx buffer with result
	return data                         // Return data read from slave register
}

func (l *LSM9DS1) I2CreadBytes(address uint, subAddress uint, dest *uint, count uint) uint {
	// var retVal []byte{}
	// l.settings.device.i2c.beginTransmission(address) // Initialize the Tx buffer
	// Next send the register to be read. OR with 0x80 to indicate multi-read.
	// l.settings.device.i2c.write(subAddress | 0x80)        // Put slave register address in Tx buffer
	// retVal = l.settings.device.i2c.endTransmission(false) // Send Tx buffer, send a restart to keep connection alive
	// endTransmission should return 0 on success
	// if retVal != 0 {
	// 	return 0
	// }
	l.settings.device.i2c.WriteByte(address, byte{byte(subAddress) | 0x80})
	// retVal = l.settings.device.i2c.requestFrom(address, count) // Read bytes from slave register address
	retVal, err := l.settings.device.i2c.ReadBytes(address, count) // Read bytes from slave register address
	if err != nil {
		log.Panic(err)
	}
	if len(retVal) != count {
		log.Panic("Expected %v recieved %v", count, len(retVal))
	}
	// for i := 0; i < count; i++ {
	// 	dest[i] = l.settings.device.i2c.read()
	// }
	// return count
	return retVal
}
