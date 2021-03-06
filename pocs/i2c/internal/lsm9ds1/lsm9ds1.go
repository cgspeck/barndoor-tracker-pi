package lsm9ds1

import (
	"fmt"
	"log"
)

func New(i2c I2CBus) (*LSM9DS1, error) {
	l := LSM9DS1{
		agAddress: byte(0x6b),
		mAddress:  byte(0x1e),
		i2c:       i2c,
		A:         &MutexReading{},
		G:         &MutexReading{},
		M:         &MutexReading{},
		magMax:    &MutexReading{},
		magMin:    &MutexReading{},
	}
	l.setDefaults()
	l.setSensitivities()

	err := l.checkWhoAmI()
	if err != nil {
		return &l, err
	}

	// Gyro initialization stuff:
	l.initGyro() // This will "turn on" the gyro. Setting up interrupts, etc.

	// // Accelerometer initialization stuff:
	l.initAccel() // "Turn on" all axes of the accel. Set up interrupts, etc.

	// // Magnetometer initialization stuff:
	l.initMag() // "Turn on" all axes of the mag. Set up interrupts, etc.

	return &l, nil
}

func (l *LSM9DS1) checkWhoAmI() error {
	mTest, err := l.mReadByteFromReg(WHO_AM_I_M)
	if err != nil {
		log.Println(err)
		return err
	}

	agTest, err := l.agReadByteFromReg(WHO_AM_I_XG)
	if err != nil {
		log.Println(err)
		return err
	}

	if mTest != WHO_AM_I_M_RSP {
		log.Println("Magnetometer whoam failed!")
		return MagnetoWhoamiFailed{}
	}

	if agTest != WHO_AM_I_AG_RSP {
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
	l.settings.gyro.scale = G_SCALE_245DPS
	// gyro sample rate: value between 1-6
	// 1 = 14.9    4 = 238
	// 2 = 59.5    5 = 476
	// 3 = 119     6 = 952
	l.settings.gyro.sampleRate = G_ODR_952
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
	l.settings.accel.scale = A_SCALE_2G
	// accel sample rate can be 1-6
	// 1 = 10 Hz    4 = 238 Hz
	// 2 = 50 Hz    5 = 476 Hz
	// 3 = 119 Hz   6 = 952 Hz
	l.settings.accel.sampleRate = XL_ODR_952
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
	l.settings.mag.scale = M_SCALE_4GS
	// mag data rate can be 0-7
	// 0 = 0.625 Hz  4 = 10 Hz
	// 1 = 1.25 Hz   5 = 20 Hz
	// 2 = 2.5 Hz    6 = 40 Hz
	// 3 = 5 Hz      7 = 80 Hz
	l.settings.mag.sampleRate = M_ODR_80
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

func (l *LSM9DS1) initGyro() {
	var tempRegValue byte = 0

	// CTRL_REG1_G (Default value: 0x00)
	// [ODR_G2][ODR_G1][ODR_G0][FS_G1][FS_G0][0][BW_G1][BW_G0]
	// ODR_G[2:0] - Output data rate selection
	// FS_G[1:0] - Gyroscope full-scale selection
	// BW_G[1:0] - Gyroscope bandwidth selection

	// To disable gyro, set sample rate bits to 0. We'll only set sample
	// rate if the gyro is enabled.
	if l.settings.gyro.enabled {
		tempRegValue = (byte(l.settings.gyro.sampleRate) & 0x07) << 5
	}

	switch l.settings.gyro.scale {
	case G_SCALE_500DPS:
		tempRegValue |= (0x1 << 3)
	case G_SCALE_2000DPS:
		tempRegValue |= (0x3 << 3)
		// Otherwise we'll set it to 245 dps (0x0 << 4)
	}

	tempRegValue |= (byte(l.settings.gyro.bandwidth) & 0x3)
	// spew.Dump(tempRegValue)
	l.agWriteToReg(CTRL_REG1_G, []byte{tempRegValue})

	// CTRL_REG2_G (Default value: 0x00)
	// [0][0][0][0][INT_SEL1][INT_SEL0][OUT_SEL1][OUT_SEL0]
	// INT_SEL[1:0] - INT selection configuration
	// OUT_SEL[1:0] - Out selection configuration
	l.agWriteToReg(CTRL_REG2_G, []byte{0x00})

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
	l.agWriteToReg(CTRL_REG3_G, []byte{tempRegValue})

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
	l.agWriteToReg(CTRL_REG4, []byte{tempRegValue})

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
	l.agWriteToReg(ORIENT_CFG_G, []byte{tempRegValue})
}

func (l *LSM9DS1) initAccel() {
	var tempRegValue byte = 0

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

	l.agWriteToReg(CTRL_REG5_XL, []byte{tempRegValue})

	// CTRL_REG6_XL (0x20) (Default value: 0x00)
	// [ODR_XL2][ODR_XL1][ODR_XL0][FS1_XL][FS0_XL][BW_SCAL_ODR][BW_XL1][BW_XL0]
	// ODR_XL[2:0] - Output data rate & power mode selection
	// FS_XL[1:0] - Full-scale selection
	// BW_SCAL_ODR - Bandwidth selection
	// BW_XL[1:0] - Anti-aliasing filter bandwidth selection
	tempRegValue = 0
	// To disable the accel, set the sampleRate bits to 0.
	if l.settings.accel.enabled {
		tempRegValue |= (byte(l.settings.accel.sampleRate) & 0x07) << 5
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
		tempRegValue |= byte(l.settings.accel.bandwidth) & 0x03
	}
	l.agWriteToReg(CTRL_REG6_XL, []byte{tempRegValue})

	// CTRL_REG7_XL (0x21) (Default value: 0x00)
	// [HR][DCF1][DCF0][0][0][FDS][0][HPIS1]
	// HR - High resolution mode (0: disable, 1: enable)
	// DCF[1:0] - Digital filter cutoff frequency
	// FDS - Filtered data selection
	// HPIS1 - HPF enabled for interrupt function
	tempRegValue = 0
	if l.settings.accel.highResEnable {
		tempRegValue |= (1 << 7) // Set HR bit
		tempRegValue |= (byte(l.settings.accel.highResBandwidth) & 0x3) << 5
	}
	l.agWriteToReg(CTRL_REG7_XL, []byte{tempRegValue})
}

func (l *LSM9DS1) initMag() {
	var tempRegValue byte = 0

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
	tempRegValue |= (byte(l.settings.mag.sampleRate) & 0x7) << 2
	l.mWriteToReg(CTRL_REG1_M, []byte{tempRegValue})

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
	l.mWriteToReg(CTRL_REG2_M, []byte{tempRegValue}) // +/-4Gauss

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
	l.mWriteToReg(CTRL_REG3_M, []byte{tempRegValue}) // Continuous conversion mode

	// CTRL_REG4_M (Default value: 0x00)
	// [0][0][0][0][OMZ1][OMZ0][BLE][0]
	// OMZ[1:0] - Z-axis operative mode selection
	//	00:low-power mode, 01:medium performance
	//	10:high performance, 10:ultra-high performance
	// BLE - Big/little endian data
	tempRegValue = 0
	tempRegValue = (l.settings.mag.ZPerformance & 0x3) << 2
	l.mWriteToReg(CTRL_REG4_M, []byte{tempRegValue})

	// CTRL_REG5_M (Default value: 0x00)
	// [0][BDU][0][0][0][0][0][0]
	// BDU - Block data update for magnetic data
	//	0:continuous, 1:not updated until MSB/LSB are read
	tempRegValue = 0
	l.mWriteToReg(CTRL_REG5_M, []byte{tempRegValue})
}

func (l *LSM9DS1) GyroAvailable() bool {
	status, err := l.agReadByteFromReg(STATUS_REG_1)
	if err != nil {
		log.Printf("error reading gyro avaliable register: %v", err)
		return false
	}

	return gyroAvailable(status)
}

func (l *LSM9DS1) AccelAvailable() bool {
	status, err := l.agReadByteFromReg(STATUS_REG_1)

	if err != nil {
		log.Printf("error reading gyro avaliable register: %v", err)
		return false
	}
	return accelAvailable(status)
}

func (l *LSM9DS1) MagAvailable(axis Axis) bool {
	status, err := l.mReadByteFromReg(STATUS_REG_M)
	if err != nil {
		log.Printf("error reading magneto avaliable register: %v", err)
		return false
	}

	return magAvailable(axis, status)
}

/*
 n.b. not supported by chip/board, see https://github.com/sparkfun/SparkFun_LSM9DS1_Arduino_Library/issues/20

func (l *LSM9DS1) TempAvailable() bool {
	status, err := l.agReadByteFromReg(STATUS_REG_1)

	if err != nil {
		log.Printf("error reading temp avaliable register: %v", err)
		return false
	}

	return ((status & (1 << 2)) >> 2) == 1
}
*/

// ReadGyro reads the Gyroscope and stores values in Gx, Gy, Gz
func (l *LSM9DS1) ReadGyro() error {
	// Read 6 bytes, beginning at OUT_X_L_G
	var raw = make([]byte, 6)
	err := l.agReadFromReg(OUT_X_L_G, raw)
	if err != nil {
		log.Printf("error reading gyro values: %v", err)
		return err
	}
	x := int16(raw[1])<<8 | int16(raw[0]) // Store x-axis values into gx
	y := int16(raw[3])<<8 | int16(raw[2]) // Store y-axis values into gy
	z := int16(raw[5])<<8 | int16(raw[4]) // Store z-axis values into gz
	if l._autoCalc {
		x -= l.gBiasRaw[X_AXIS]
		y -= l.gBiasRaw[Y_AXIS]
		z -= l.gBiasRaw[Z_AXIS]
	}
	l.G.SetReading(x, y, z)
	return nil
}

func (l *LSM9DS1) ReadAccel() error {
	// We'll read six bytes from the accelerometer into temp
	var raw = make([]byte, 6)

	// Read 6 bytes, beginning at OUT_X_L_XL

	err := l.agReadFromReg(OUT_X_L_XL, raw)
	if err != nil {
		log.Printf("error reading accelerometer values: %v", err)
		return err
	}
	x := int16(raw[1])<<8 | int16(raw[0]) // Store x-axis values into ax
	y := int16(raw[3])<<8 | int16(raw[2]) // Store y-axis values into ay
	z := int16(raw[5])<<8 | int16(raw[4]) // Store z-axis values into az
	if l._autoCalc {
		x -= l.aBiasRaw[X_AXIS]
		y -= l.aBiasRaw[Y_AXIS]
		z -= l.aBiasRaw[Z_AXIS]
	}
	l.A.SetReading(x, y, z)
	return nil
}

func (l *LSM9DS1) ReadMag() error {
	var raw = make([]byte, 6)
	// We'll read six bytes from the mag into temp

	err := l.mReadFromReg(OUT_X_L_M, raw)
	if err != nil {
		log.Printf("error reading magnetometer values: %v", err)
		return err
	}
	//
	//  first attempt to decode the 2's complement based on c code
	//
	// Read 6 bytes, beginning at OUT_X_L_M
	// l.Mx = (raw[1] << 8) | raw[0] // Store x-axis values into mx
	// l.My = (raw[3] << 8) | raw[2] // Store y-axis values into my
	// l.Mz = (raw[5] << 8) | raw[4] // Store z-axis values into mz
	//
	//  try to follow example from embd library for other magnetometer sensor:
	//
	// https://github.com/kidoman/embd/blob/master/sensor/lsm303/lsm303.go#L102

	x := int16(raw[1])<<8 | int16(raw[0]) // Store x-axis values into mx
	y := int16(raw[3])<<8 | int16(raw[2]) // Store y-axis values into my
	z := int16(raw[5])<<8 | int16(raw[4]) // Store z-axis values into mz

	l.M.SetReading(x, y, z)
	return nil
}

/*
//  n.b. not supported by chip/board? see https://github.com/sparkfun/SparkFun_LSM9DS1_Arduino_Library/issues/20

func (l *LSM9DS1) ReadTemp() {
	var raw = make([]byte, 2)
	err := l.mReadFromReg(OUT_TEMP_L, raw)
	if err != nil {
		log.Printf("Error reading temp value: %v", err)
	}
	var offset int16 = 25 // Per datasheet sensor outputs 0 typically @ 25 degrees centigrade
	l.Temperature = offset + (int16(raw[1])<<8 | int16(raw[0]))
}
*/

// Calibrate uses the FIFO to accumulate sample of accelerometer and gyro data, average
// them, scales them to  gs and deg/s, respectively, and then passes the biases to the main sketch
// for subtraction from all subsequent data. There are no gyro and accelerometer bias registers to store
// the data as there are in the ADXL345, a precursor to the LSM9DS0, or the MPU-9150, so we have to
// subtract the biases ourselves. This results in a more accurate measurement in general and can
// remove errors due to imprecise or varying initial placement. Calibration of sensor data in this manner
// is good practice.
func (l *LSM9DS1) Calibrate(autoCalc bool) {
	var samples int16 = 0
	var i int16

	aBiasRawTemp := [3]int16{0, 0, 0}
	gBiasRawTemp := [3]int16{0, 0, 0}

	// Turn on FIFO and set threshold to 32 samples
	l.enableFIFO(true)
	l.setFIFO(FIFO_THS, 0x1F)
	for samples < 31 {
		// samples = (l.xgReadByte(FIFO_SRC) & 0x3F) // Read number of stored samples
		l.agReadByteFromReg(FIFO_SRC)
		samples += 1
	}
	for i = 0; i < samples; i++ {
		// Read the gyro data stored in the FIFO
		l.ReadGyro()
		gx, gy, gz := l.G.GetReading()
		gBiasRawTemp[0] += gx
		gBiasRawTemp[1] += gy
		gBiasRawTemp[2] += gz
		l.ReadAccel()
		ax, ay, az := l.A.GetReading()
		aBiasRawTemp[0] += ax
		aBiasRawTemp[1] += ay
		aBiasRawTemp[2] += az - int16(1./l.aRes) // Assumes sensor facing up!
	}

	for i = 0; i < 3; i++ {
		l.gBiasRaw[i] = gBiasRawTemp[i] / samples
		l.gBias[i] = l.CalcGyro(l.gBiasRaw[i])
		l.aBiasRaw[i] = aBiasRawTemp[i] / samples
		l.aBias[i] = l.CalcAccel(l.aBiasRaw[i])
	}

	l.enableFIFO(false)
	l.setFIFO(FIFO_OFF, 0x00)

	fmt.Printf("Calibration gBias: %v\n", l.gBias)
	fmt.Printf("Calibration aBias: %v\n", l.aBias)

	if autoCalc {
		fmt.Print("Calibration will be accounted for in calculations\n")
		l._autoCalc = true
	}
}

// CalibrateMag takes readings from the Magnetometer and optionally
// loads the calculated bias into the chip.
func (l *LSM9DS1) CalibrateMag() {
	var i, j int
	magMin := l.magMin.ToList()
	magMax := l.magMax.ToList() // The road warrior

	for i = 0; i < 128; i++ {
		for l.MagAvailable(ALL_AXIS) == false {
		}
		l.ReadMag()
		mx, my, mz := l.M.GetReading()

		magTemp := [3]int16{0, 0, 0}
		magTemp[0] = mx
		magTemp[1] = my
		magTemp[2] = mz
		for j = 0; j < 3; j++ {
			if magTemp[j] > magMax[j] {
				magMax[j] = magTemp[j]
			}
			if magTemp[j] < magMin[j] {
				magMin[j] = magTemp[j]
			}
		}
	}
	l.magMin.FromList(magMin)
	l.magMax.FromList(magMax)
}

// Returns difference between min and max values
func (l *LSM9DS1) MagRange() []int16 {
	magMin := l.magMin.ToList()
	magMax := l.magMax.ToList()
	diff := []int16{
		magMax[0] - magMin[0],
		magMax[1] - magMin[1],
		magMax[2] - magMin[2],
	}

	return diff
}

func (l *LSM9DS1) LoadMagBias() error {
	var j int
	magMin := l.magMin.ToList()
	magMax := l.magMax.ToList() // The road warrior

	for j = 0; j < 3; j++ {
		l.mBiasRaw[j] = (magMax[j] + magMin[j]) / 2
		l.mBias[j] = l.CalcMag(l.mBiasRaw[j])
		err := l.magOffset(Axis(j), l.mBiasRaw[j])
		if err != nil {
			return err
		}
	}
	log.Printf("CalibrateMag mBias: %v", l.mBias)
	return nil
}

func (l *LSM9DS1) CalcMag(mag int16) float32 {
	// Return the mag raw reading times our pre-calculated Gs / (ADC tick):
	return l.mRes * float32(mag)
}

func (l *LSM9DS1) magOffset(axis Axis, offset int16) error {
	if axis > 2 {
		return nil
	}
	bAxis := byte(axis)
	log.Printf("magOffset axis: %v bias: %v", axis, offset)
	return l.mWriteWordToReg(OFFSET_X_REG_L_M+2*bAxis, uint16(offset))
}

func (l *LSM9DS1) enableFIFO(enable bool) error {
	temp, err := l.agReadByteFromReg(CTRL_REG9)
	if err != nil {
		log.Printf("Error reading CTRL_REG9 while enabling/disabling FIFO! desired state: %v, error: %v\n", enable, err)
		return err
	}

	iTemp := int(temp)

	if enable {
		iTemp |= (1 << 1)
	} else {
		iTemp &= ^(1 << 1)
	}

	var res byte = 0
	if iTemp < 0 {
		res = 1 - byte(iTemp)
	} else {
		res = byte(iTemp)
	}

	err = l.agWriteToReg(CTRL_REG9, []byte{res})
	if err != nil {
		log.Printf("Error writing CTRL_REG9 while enabling/disabling FIFO! desired state: %v, error: %v\n", enable, err)
		return err
	}
	return err
}

func (l *LSM9DS1) setFIFO(fifoMode FIFOMode, fifoThs byte) {
	// Limit threshold - 0x1F (31) is the maximum. If more than that was asked
	// limit it to the maximum.
	var threshold byte
	if fifoThs <= 0x1F {
		threshold = fifoThs
	} else {
		threshold = 0x1F
	}
	val := ((byte(fifoMode) & 0x7) << 5) | (threshold & 0x1F)
	l.agWriteToReg(FIFO_CTRL, []byte{val})
}

func (l *LSM9DS1) CalcGyro(gyro int16) float32 {
	// Return the gyro raw reading times our pre-calculated DPS / (ADC tick):
	return l.gRes * float32(gyro)
}

func (l *LSM9DS1) CalcAccel(accel int16) float32 {
	// Return the accel raw reading times our pre-calculated g's / (ADC tick):
	return l.aRes * float32(accel)
}

// setSensitivities sets the sensitivity (also referred to as resolution) for each sensor
func (l *LSM9DS1) setSensitivities() {
	l.aRes = aScaleSensitivity[l.settings.accel.scale]
	l.gRes = gScaleSensitivity[l.settings.gyro.scale]
	l.mRes = mScaleSensitivity[l.settings.mag.scale]
}

// IO

// reads len(dest) bytes from register
func (l *LSM9DS1) agReadFromReg(regAddress byte, dest []byte) error {
	return l.i2c.ReadFromReg(l.agAddress, regAddress, dest)
}

// reads a byte from register
func (l *LSM9DS1) agReadByteFromReg(regAddress byte) (byte, error) {
	return l.i2c.ReadByteFromReg(l.agAddress, regAddress)
}

// func (l *LSM9DS1) mReadByteFromReg(regAddress byte) (byte, error) {
// 	return l.i2c.ReadByteFromReg(l.mAddress, regAddress)
// }

// func (l *LSM9DS1) agReadBytes(regAddress byte, count int) ([]byte, error) {
// 	// have to write to device the register we want to read from
// 	// 0x80 indicates a "multiread"
// 	l.i2c.WriteByte(l.agAddress, regAddress|0x80)
// 	return l.i2c.ReadBytes(l.agAddress, count)
// }

func (l *LSM9DS1) mReadFromReg(regAddress byte, dest []byte) error {
	return l.i2c.ReadFromReg(l.mAddress, regAddress, dest)
}

func (l *LSM9DS1) mReadByteFromReg(regAddress byte) (byte, error) {
	return l.i2c.ReadByteFromReg(l.mAddress, regAddress)
}

func (l *LSM9DS1) agWriteToReg(regAddress byte, data []byte) error {
	return l.i2c.WriteToReg(l.agAddress, regAddress, data)
}

func (l *LSM9DS1) mWriteToReg(regAddress byte, data []byte) error {
	return l.i2c.WriteToReg(l.mAddress, regAddress, data)
}

func (l *LSM9DS1) mWriteWordToReg(regAddress byte, data uint16) error {
	return l.i2c.WriteWordToReg(l.mAddress, regAddress, data)
}

func (l *LSM9DS1) agWriteWordToReg(regAddress byte, data uint16) error {
	return l.i2c.WriteWordToReg(l.agAddress, regAddress, data)
}
