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

			input := make(map[string]interface{})
			err = json.Unmarshal(body, &input)
			if err != nil {
				return 500, err
			}
			err = ah.SetLocationSettings(input)
			if err != nil {
				return 500, err
			}
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
