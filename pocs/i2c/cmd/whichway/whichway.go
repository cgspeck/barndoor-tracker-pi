package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/cgspeck/bdt/pocs/i2c/internal/lsm9ds1"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

func main() {
	bus := embd.NewI2CBus(1)
	defer bus.Close()

	l, err := lsm9ds1.New(bus)
	if err != nil {
		fmt.Printf("Error instantiating driver: %v", err)
		os.Exit(1)
	}

	fmt.Println("Begin calibration")
	l.Calibrate(true)
	fmt.Println("End calibration")

	fmt.Println("Begin Magneto calibration")
	l.CalibrateMag(true)
	fmt.Println("End Magneto calibration")

	printInterval := time.Millisecond * 250
	declination := 11.72
	lastPrint := time.Now()

	for true {
		current := time.Now()
		if current.Sub(lastPrint) >= printInterval {
			if l.GyroAvailable() {
				l.ReadGyro()
				fmt.Printf("Gyro read (deg/s): x=%v y=%v z=%v\n",
					l.CalcGyro(l.Gx),
					l.CalcGyro(l.Gy),
					l.CalcGyro(l.Gz),
				)
			}

			if l.AccelAvailable() {
				l.ReadAccel()
				fmt.Printf("Accel read (g's): x=%v y=%v z=%v\n",
					l.CalcAccel(l.Ax),
					l.CalcAccel(l.Ay),
					l.CalcAccel(l.Az),
				)
			}

			if l.MagAvailable(lsm9ds1.ALL_AXIS) {
				l.ReadMag()
				fmt.Printf("Magneto read (gauss): x=%v y=%v z=%v\n",
					l.CalcMag(l.Mx),
					l.CalcMag(l.My),
					l.CalcMag(l.Mz),
				)
			}

			printAttitude(l.Ax, l.Ay, l.Az, -l.My, -l.Mx, l.Mz, declination)
			lastPrint = current
		}
	}
}

func printAttitude(ax int16, ay int16, az int16, mx int16, my int16, mz int16, declination float64) {
	fax := float64(ax)
	fay := float64(ay)
	faz := float64(az)

	roll := math.Atan2(fay, faz)
	pitch := math.Atan2(-fax, math.Sqrt(fay*fay+faz*faz))

	heading := math.Atan2(float64(my), float64(mx)) / math.Pi * 180
	if heading < 0 {
		heading += 360
	}
	
	// Convert everything from radians to degrees:
	pitch *= 180.0 / math.Pi
	roll *= 180.0 / math.Pi

	fmt.Printf("\n\nPitch, Roll: %v, %v				Heading: %v\n\n", pitch, roll, heading)
}
