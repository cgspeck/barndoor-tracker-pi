package ds18b20_wrapper

type WrappedDS18B20 struct{}

func New() (*WrappedDS18B20, error) {

	return &WrappedDS18B20{}, nil
}

func (_ *WrappedDS18B20) Temperature() float64 {
	return 0
}

func (_ *WrappedDS18B20) SensorOk() bool {
	return false
}
