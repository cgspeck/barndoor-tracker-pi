package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cgspeck/bdt/pocs/i2c/internal/lsm9ds1"

	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

func main() {
	bus := lsm9ds1.NewMutexI2cBus(1)
	defer bus.Close()

	l, err := lsm9ds1.New(&bus)
	if err != nil {
		fmt.Printf("Error instantiating driver: %v", err)
		os.Exit(1)
	}

	// the next two lines can be called repeatedly until calibration looks good
	l.CalibrateMag()
	fmt.Printf("Mag range: %v\n", l.MagRange)
	l.LoadMagBias()
	printInterval := time.Millisecond * 250
	lastPrint := time.Now()

	for true {
		current := time.Now()
		if current.Sub(lastPrint) >= printInterval {
			if l.MagAvailable(lsm9ds1.ALL_AXIS) {
				l.ReadMag()
				mx, my, mz := l.M.GetReading()
				fmt.Printf("Magneto calc (gauss): x=%v y=%v z=%v\n",
					l.CalcMag(mx),
					l.CalcMag(my),
					l.CalcMag(mz),
				)
			}

			lastPrint = current
		}
	}
}
