package ds18b20_wrapper

type IWrappedDS18B20 interface {
	Temperature() float64
	SensorOk() bool
}
