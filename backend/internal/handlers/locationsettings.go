package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type LocationManagementDisabledError struct{}

func (LocationManagementDisabledError) Error() string {
	return "Location management is disabled"
}

func LocationSettingsHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "POST" {
		locationSettings := ah.GetLocationSettings()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return 500, err
		}

		input := make(map[string]interface{})
		err = json.Unmarshal(body, &input)
		if err != nil {
			return 500, err
		}

		locationSettings.Lock()
		err = ah.SetLocationSettings(input)
		locationSettings.Unlock()
		if err != nil {
			return 500, err
		}
	}

	if r.Method == "GET" || r.Method == "POST" {
		locationSettings := ah.GetLocationSettings()
		locationSettings.RLock()
		defer locationSettings.RUnlock()
		err := writeJson(locationSettings, w)
		if err != nil {
			return 500, err
		}

		return 200, nil
	}

	return 501, nil
}
