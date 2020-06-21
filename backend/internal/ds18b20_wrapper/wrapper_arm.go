package ds18b20_wrapper

import (
	"log"

	"github.com/yryz/ds18b20"
)

type WrappedDS18B20 struct {
	sensorOk      bool
	sensorAddress string
}

func New() (*WrappedDS18B20, error) {
	sensors, err := ds18b20.Sensors()
	sensorOk := false
	var sensorAddress string

	if err != nil {
		log.Printf("Error initialising ds18b20 library: %v\n", err)
	}

	sensorCount := len(sensors)
	log.Printf("Found %v sensors\n", sensorCount)

	if sensorCount > 0 {
		log.Printf("All sensors: %v\n", sensors)
		sensorAddress = sensors[0]
		log.Printf("Selected sensor %q\n", sensorAddress)
	}

	return &WrappedDS18B20{
		sensorOk:      sensorOk,
		sensorAddress: sensorAddress,
	}, nil
}

func (w *WrappedDS18B20) Temperature() float64 {
	v, err := ds18b20.Temperature(w.sensorAddress)

	if err != nil {
		log.Printf("Error reading ds18b20 sensor %q: %v\n", w.sensorAddress, err)
	}

	return v
}

func (w *WrappedDS18B20) SensorOk() bool {
	return w.sensorOk
}
