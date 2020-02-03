package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type LocationManagementDisabledError struct{}

func (LocationManagementDisabledError) Error() string {
	return "Location management is disabled"
}

func LocationSettingsHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "GET" {
		err := writeJson(ah.GetLocationSettings(), w)
		if err != nil {
			return 500, err
		}

		return 200, nil
	}

	if r.Method == "POST" {
		if ah.GetLocationSettings().ManagementEnabled {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return 500, err
			}
			var LocationSettings models.LocationSettings
			err = json.Unmarshal(body, &LocationSettings)
			if err != nil {
				return 500, err
			}

			// if ah.GetLocationSettings().AccessPointMode != LocationSettings.AccessPointMode {
			// 	ah.SetAPMode(LocationSettings.AccessPointMode)
			// }
			err = writeJson(ah.GetLocationSettings(), w)
			if err != nil {
				return 500, err
			}
		} else {
			return 400, LocationManagementDisabledError{}
		}
	}

	return 501, nil
}
