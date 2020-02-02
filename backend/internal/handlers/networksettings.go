package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func NetworkSettingsHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "GET" || r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return 500, err
			}
			var networkSettings models.NetworkSettings
			err = json.Unmarshal(body, &networkSettings)
			if err != nil {
				return 500, err
			}

			if ah.GetNetworkSettings().AccessPointMode != networkSettings.AccessPointMode {
				ah.SetAPMode(networkSettings.AccessPointMode)
			}
		}

		b, err := json.MarshalIndent(ah.GetNetworkSettings(), "", "  ")
		if err != nil {
			return 500, err
		}
		io.WriteString(w, string(b))

		return 200, nil
	}
	return 501, nil
}
