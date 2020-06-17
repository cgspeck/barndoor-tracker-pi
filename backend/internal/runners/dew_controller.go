package runners

type DewControllerRunner struct{}

func NewDewControllerRunner() (*DewControllerRunner, error) {
	return &DewControllerRunner{}, nil
}
