package handlers

import (
	"net/http"
)

func AlignHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	alignStatus := ah.GetAlignStatus()
	alignStatus.RLock()
	defer alignStatus.RUnlock()

	err := writeJson(alignStatus, w)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
