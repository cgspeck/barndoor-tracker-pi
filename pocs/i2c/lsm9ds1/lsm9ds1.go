package lsm9ds1

import (
	"log"

	"github.com/kidoman/embd"
)

func New(i2c embd.I2CBus) (*LSM9DS1, error) {
	l := LSM9DS1{
		agAddress: byte(0x6b),
		mAddress:  byte(0x1e),
		i2c:       i2c,
	}
	l.setDefaults()
	err := l.checkWhoAmI()

	return &l, err
}

func (l *LSM9DS1) checkWhoAmI() error {
	mTest := make([]byte, 1)
	agTest := make([]byte, 1)

	err := l.mReadReg(WHO_AM_I_M, mTest)
	if err != nil {
		log.Println(err)
		return err
	}

	err = l.agReadReg(WHO_AM_I_XG, agTest)
	if err != nil {
		log.Println(err)
		return err
	}

	if mTest[0] != WHO_AM_I_M_RSP {
		log.Println("Magnetometer whoam failed!")
		return MagnetoWhoamiFailed{}
	}

	if agTest[0] != WHO_AM_I_AG_RSP {
		log.Println("Accel/Gyro whoam failed!")
		return AGWhoamiFailed{}
	}
	log.Println("whoami check pass")
	return nil
}

func (l *LSM9DS1) setDefaults() {
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

// IO
func (l *LSM9DS1) agReadReg(regAddress byte, dest []byte) error {
	return l.i2c.ReadFromReg(l.agAddress, regAddress, dest)
}

func (l *LSM9DS1) mReadReg(regAddress byte, dest []byte) error {
	return l.i2c.ReadFromReg(l.mAddress, regAddress, dest)
}
