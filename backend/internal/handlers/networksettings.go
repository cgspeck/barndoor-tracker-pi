package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type NetworkManagementDisabledError struct{}

func (NetworkManagementDisabledError) Error() string {
	return "Network management is disabled"
}

func NetworkSettingsHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "GET" {
		err := writeJson(ah.GetNetworkSettings(), w)
		if err != nil {
			return 500, err
		}

		return 200, nil
	}

	if r.Method == "POST" {
		if ah.GetNetworkSettings().ManagementEnabled {
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
			err = writeJson(ah.GetNetworkSettings(), w)
			if err != nil {
				return 500, err
			}
		} else {
			return 400, NetworkManagementDisabledError{}
		}
	}

	return 501, nil
}
