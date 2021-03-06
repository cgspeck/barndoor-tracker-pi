package handlers

import (
	"net/http"
)

func PIDLogHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	logs := ah.GetPIDLogFiles()

	err := writeJson(logs, w)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
