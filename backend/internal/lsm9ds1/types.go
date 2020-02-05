package lsm9ds1

import "github.com/kidoman/embd"

/******************************************************************************
SFE_LSM9DS1 Library - LSM9DS1 Types and Enumerations
Ported from library by [SparkFun Electronics](https://github.com/sparkfun/LSM9DS1_Breakout)

This file defines all types and enumerations used by the LSM9DS1 struct.

******************************************************************************/

// The LSM9DS1 functions over both I2C or SPI. This library supports both.
// But the interface mode used must be sent to the LSM9DS1 constructor. Use
// one of these two as the first parameter of the constructor.
const (
	IMU_MODE_SPI = iota
	IMU_MODE_I2C
)

// accel_scale defines all possible FSR's of the accelerometer:
const (
	A_SCALE_2G  = iota // 00:  2g
	A_SCALE_16G        // 01:  16g
	A_SCALE_4G         // 10:  4g
	A_SCALE_8G         // 11:  8g
)

// gyro_scale defines the possible full-scale ranges of the gyroscope:
const (
	G_SCALE_245DPS  = iota // 00:  245 degrees per second
	G_SCALE_500DPS         // 01:  500 dps
	G_SCALE_2000DPS        // 11:  2000 dps
)

// mag_scale defines all possible FSR's of the magnetometer:
const (
	M_SCALE_4GS  = iota // 00:  4Gs
	M_SCALE_8GS         // 01:  8Gs
	M_SCALE_12GS        // 10:  12Gs
	M_SCALE_16GS        // 11:  16Gs
)

// gyro_odr defines all possible data rate/bandwidth combos of the gyro:
const (
	//! TODO
	G_ODR_PD  = iota // Power down (0)
	G_ODR_149        // 14.9 Hz (1)
	G_ODR_595        // 59.5 Hz (2)
	G_ODR_119        // 119 Hz (3)
	G_ODR_238        // 238 Hz (4)
	G_ODR_476        // 476 Hz (5)
	G_ODR_952        // 952 Hz (6)
)

// accel_oder defines all possible output data rates of the accelerometer:
const (
	XL_POWER_DOWN = iota // Power-down mode (0x0)
	XL_ODR_10            // 10 Hz (0x1)
	XL_ODR_50            // 50 Hz (0x02)
	XL_ODR_119           // 119 Hz (0x3)
	XL_ODR_238           // 238 Hz (0x4)
	XL_ODR_476           // 476 Hz (0x5)
	XL_ODR_952           // 952 Hz (0x6)
)

// accel_abw defines all possible anti-aliasing filter rates of the accelerometer:
const (
	A_ABW_408 = iota // 408 Hz (0x0)
	A_ABW_211        // 211 Hz (0x1)
	A_ABW_105        // 105 Hz (0x2)
	A_ABW_50         //  50 Hz (0x3)
)

// mag_odr defines all possible output data rates of the magnetometer:
const (
	M_ODR_0625 = iota // 0.625 Hz (0)
	M_ODR_125         // 1.25 Hz (1)
	M_ODR_250         // 2.5 Hz (2)
	M_ODR_5           // 5 Hz (3)
	M_ODR_10          // 10 Hz (4)
	M_ODR_20          // 20 Hz (5)
	M_ODR_40          // 40 Hz (6)
	M_ODR_80          // 80 Hz (7)
)

const (
	XG_INT1 = 0x0C // INT1_CTRL
	XG_INT2 = 0x0D // INT2_CTRL
)

const (
	INT_DRDY_XL    = (1 << 0) // Accelerometer data ready (INT1 & INT2)
	INT_DRDY_G     = (1 << 1) // Gyroscope data ready (INT1 & INT2)
	INT1_BOOT      = (1 << 2) // Boot status (INT1)
	INT2_DRDY_TEMP = (1 << 2) // Temp data ready (INT2)
	INT_FTH        = (1 << 3) // FIFO threshold interrupt (INT1 & INT2)
	INT_OVR        = (1 << 4) // Overrun interrupt (INT1 & INT2)
	INT_FSS5       = (1 << 5) // FSS5 interrupt (INT1 & INT2)
	INT_IG_XL      = (1 << 6) // Accel interrupt generator (INT1)
	INT1_IG_G      = (1 << 7) // Gyro interrupt enable (INT1)
	INT2_INACT     = (1 << 7) // Inactivity interrupt output (INT2)
)

const (
	XLIE_XL = (1 << 0)
	XHIE_XL = (1 << 1)
	YLIE_XL = (1 << 2)
	YHIE_XL = (1 << 3)
	ZLIE_XL = (1 << 4)
	ZHIE_XL = (1 << 5)
	GEN_6D  = (1 << 6)
)

const (
	XLIE_G = (1 << 0)
	XHIE_G = (1 << 1)
	YLIE_G = (1 << 2)
	YHIE_G = (1 << 3)
	ZLIE_G = (1 << 4)
	ZHIE_G = (1 << 5)
)

const (
	ZIEN = (1 << 5)
	YIEN = (1 << 6)
	XIEN = (1 << 7)
)

const (
	INT_ACTIVE_HIGH = iota
	INT_ACTIVE_LOW
)

const (
	INT_PUSH_PULL = iota
	INT_OPEN_DRAIN
)

const (
	FIFO_OFF          = 0
	FIFO_THS          = 1
	FIFO_CONT_TRIGGER = 3
	FIFO_OFF_TRIGGER  = 4
	FIFO_CONT         = 6
)

type gyroSettings struct {
	// Gyroscope settings:
	enabled        bool
	scale          uint // Changed this to 16-bit
	sampleRate     uint
	bandwidth      uint
	lowPowerEnable uint
	HPFEnable      uint
	HPFCutoff      uint
	flipX          uint
	flipY          uint
	flipZ          uint
	orientation    uint
	enableX        bool
	enableY        bool
	enableZ        bool
	latchInterrupt uint
}

type deviceSettings struct {
	commInterface uint // Can be I2C SPI 4-wire or SPI 3-wire
	agAddress     uint // I2C address or SPI CS pin
	mAddress      uint // I2C address or SPI CS pin
	TwoWire       *embd.I2CBus
	//   TwoWire* i2c;    // pointer to an instance of I2C interface
}

type accelSettings struct {
	// Accelerometer settings:
	enabled          bool
	scale            uint
	sampleRate       uint
	enableX          bool
	enableY          bool
	enableZ          bool
	bandwidth        int
	highResEnable    uint
	highResBandwidth uint
}

type magSettings struct {
	// Magnetometer settings:
	enabled                bool
	scale                  uint
	sampleRate             uint
	tempCompensationEnable uint
	XYPerformance          uint
	ZPerformance           uint
	lowPowerEnable         bool
	operatingMode          uint
}

type temperatureSettings struct {
	// Temperature settings
	enabled bool
}

type IMUSettings struct {
	device deviceSettings
	gyro   gyroSettings
	accel  accelSettings
	mag    magSettings
	temp   temperatureSettings
}

type LSM9DS1 struct {
	settings                     IMUSettings
	gx, gy, gz                   int // x, y, and z axis readings of the gyroscope
	ax, ay, az                   int // x, y, and z axis readings of the accelerometer
	mx, my, mz                   int // x, y, and z axis readings of the magnetometer
	temperature                  int // Chip temperature
	gBias, aBias, mBias          [3]float32
	gBiasRaw, aBiasRaw, mBiasRaw [3]int

	// protected
	// x_mAddress and gAddress store the I2C address or SPI chip select pin
	// for each sensor.
	_mAddress, _xgAddress uint

	// gRes, aRes, and mRes store the current resolution for each sensor.
	// Units of these values would be DPS (or g's or Gs's) per ADC tick.
	// This value is calculated as (sensor scale) / (2^15).
	gRes, aRes, mRes float32

	// _autoCalc keeps track of whether we're automatically subtracting off
	// accelerometer and gyroscope bias calculated in calibrate().
	_autoCalc bool
}
