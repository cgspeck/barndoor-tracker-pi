package handlers

import (
	"net/http"
)

func Flags(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	flags := ah.GetFlags()
	flags.RLock()
	defer flags.RUnlock()

	err := writeJson(flags, w)
	if err != nil {
		return 500, err
	}
	return 200, nil
}
