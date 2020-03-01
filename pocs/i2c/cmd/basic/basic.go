package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cgspeck/bdt/pocs/i2c/internal/lsm9ds1"

	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

func main() {
	log.Println("hello world")
	// bus := embd.NewI2CBus(1)
	bus := lsm9ds1.NewMutexI2cBus(1)
	defer bus.Close()

	l, err := lsm9ds1.New(&bus)
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
	lastPrint := time.Now()

	for true {
		current := time.Now()
		if current.Sub(lastPrint) >= printInterval {
			if l.GyroAvailable() {
				l.ReadGyro()
				gx, gy, gz := l.G.GetReading()
				fmt.Printf("Gyro read: x=%v y=%v z=%v\n", gx, gy, gz)
			}

			if l.AccelAvailable() {
				l.ReadAccel()
				ax, ay, az := l.A.GetReading()
				fmt.Printf("Accel read: x=%v y=%v z=%v\n", ax, ay, az)
			}

			if l.MagAvailable(lsm9ds1.ALL_AXIS) {
				l.ReadMag()
				mx, my, mz := l.M.GetReading()
				fmt.Printf("Magneto read: x=%v y=%v z=%v\n", mx, my, mz)
			}

			lastPrint = current
		}
	}
}
