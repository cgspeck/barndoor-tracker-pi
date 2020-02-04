package handlers

import (
	"net/http"
)

func Flags(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Content-Type", "application/json")
		return 204, nil
	}

	w.Header().Set("Content-Type", "application/json")
	err := writeJson(ah.GetFlags(), w)
	if err != nil {
		return 500, err
	}
	return 200, nil
}
