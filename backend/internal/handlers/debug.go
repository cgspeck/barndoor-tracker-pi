package handlers

import (
	"net/http"
)

func DebugHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	err := writeJson(ah.GetContext(), w)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
