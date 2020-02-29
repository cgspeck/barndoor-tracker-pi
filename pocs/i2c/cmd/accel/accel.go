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

	l.Calibrate(true)

	printInterval := time.Millisecond * 250
	lastPrint := time.Now()

	for true {
		current := time.Now()
		if current.Sub(lastPrint) >= printInterval {
			if l.AccelAvailable() {
				l.ReadAccel()
				fmt.Printf("Accel read (g's): x=%v y=%v z=%v\n",
					l.CalcAccel(l.Ax),
					l.CalcAccel(l.Ay),
					l.CalcAccel(l.Az),
				)
			}
			lastPrint = current
		}
	}
}
