package main

import (
	"log"

	"./lsm9ds1"

	// "github.com/davecgh/go-spew/spew"

	"github.com/davecgh/go-spew/spew"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

func main() {
	log.Println("hello world")
	bus := embd.NewI2CBus(1)
	defer bus.Close()

	err, l := lsm9ds1.New(bus)
	spew.Dump(err, l)
}
