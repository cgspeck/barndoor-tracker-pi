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
	networkSettings := ah.GetNetworkSettings()

	if r.Method == "POST" {
		networkSettings.RLock()
		managementEnabled := networkSettings.ManagementEnabled
		accessPointMode := networkSettings.AccessPointMode
		networkSettings.RUnlock()

		if managementEnabled {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return 500, err
			}
			var newNetworkSettings models.NetworkSettings
			err = json.Unmarshal(body, &newNetworkSettings)
			if err != nil {
				return 500, err
			}

			if newNetworkSettings.AccessPointMode != accessPointMode {
				networkSettings.Lock()
				ah.SetAPMode(newNetworkSettings.AccessPointMode)
				networkSettings.AccessPointMode = newNetworkSettings.AccessPointMode
				networkSettings.Unlock()
			}

		} else {
			return 400, NetworkManagementDisabledError{}
		}
	}

	if r.Method == "GET" || r.Method == "POST" {
		networkSettings.RLock()
		defer networkSettings.RUnlock()
		err := writeJson(networkSettings, w)
		if err != nil {
			return 500, err
		}
		return 200, nil
	}

	return 501, nil
}
