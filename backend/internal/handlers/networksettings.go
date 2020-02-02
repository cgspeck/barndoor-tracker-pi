package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func NetworkSettingsHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "GET" {
		b, err := json.MarshalIndent(ah.GetNetworkSettings(), "", "  ")
		if err != nil {
			return 500, err
		}
		io.WriteString(w, string(b))

		return 200, nil
	}

	return 501, nil
}
