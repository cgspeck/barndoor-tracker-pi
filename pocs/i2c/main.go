package main

import (
	"log"

	"./lsm9ds1"

	// "github.com/davecgh/go-spew/spew"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

func main() {
	log.Println("hello world")
	bus := embd.NewI2CBus(1)
	defer bus.Close()

	mAddress := 0x1e
	agAddress := 0x6b
	mTest := make([]byte, 1)
	agTest := make([]byte, 1)

	err := bus.ReadFromReg(byte(mAddress), lsm9ds1.WHO_AM_I_M, mTest)
	if err != nil {
		log.Fatalln(err)
	}

	err = bus.ReadFromReg(byte(agAddress), lsm9ds1.WHO_AM_I_XG, agTest)
	if err != nil {
		log.Fatalln(err)
	}

	if mTest[0] != lsm9ds1.WHO_AM_I_M_RSP {
		log.Fatalln("Magnetometer whoam failed!")
	}

	if agTest[0] != lsm9ds1.WHO_AM_I_AG_RSP {
		log.Fatalln("Accel/Gyro whoam failed!")
	}
	log.Println("whoam check pass")
}
