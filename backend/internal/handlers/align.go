package handlers

import (
	"net/http"
)

func AlignHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	err := writeJson(ah.GetAlignStatus(), w)
	if err != nil {
		return 500, err
	}

	return 200, nil
}