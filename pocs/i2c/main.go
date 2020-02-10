package main

import (
	"fmt"
	"log"
	"os"

	"./lsm9ds1"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

func main() {
	log.Println("hello world")
	bus := embd.NewI2CBus(1)
	defer bus.Close()

	l, err := lsm9ds1.New(bus)
	if err != nil {
		fmt.Printf("Error instantiating driver: %v", err)
		os.Exit(1)
	}
	for true {
		if l.GyroAvailable() {
			l.ReadGyro()
			fmt.Printf("Gyro read: x=%v y=%v z=%v\n", l.Gx, l.Gy, l.Gz)
		}

		if l.AccelAvailable() {
			fmt.Printf("Accel read: x=%v y=%v z=%v\n", l.Ax, l.Ay, l.Az)
		}
	}
}
