package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LocationManagementDisabledError struct{}

func (LocationManagementDisabledError) Error() string {
	return "Location management is disabled"
}

type LocationManagementMissingKeyError struct {
	k string
}

func (e LocationManagementMissingKeyError) Error() string {
	return fmt.Sprintf("Key missing %q", e.k)
}

type LocationManagementInvalidValueError struct {
	v interface{}
}

func (e LocationManagementInvalidValueError) Error() string {
	return fmt.Sprintf("Invalid value %q", e.v)
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

		path := r.URL.Path

		switch path {
		case "/backend/settings/location":
			err = ah.SetLocationSettings(input)
			fmt.Println("b5")
		case "/backend/toggle/ignoreAz":
			fallthrough
		case "/backend/toggle/ignoreAlt":
			rv, ok := input["enabled"]

			if !ok {
				return 400, LocationManagementMissingKeyError{"enabled"}
			}

			bv, ok := rv.(bool)

			if !ok {
				return 400, LocationManagementInvalidValueError{rv}
			}

			locationSettings.Lock()
			if path == "/backend/toggle/ignoreAz" {
				locationSettings.IgnoreAz = bv
			} else {
				locationSettings.IgnoreAlt = bv
			}
			locationSettings.Unlock()
		}

		if err != nil {
			return 500, err
		}
	}

	if r.Method == "GET" || r.Method == "POST" {
		fmt.Println("b6")
		locationSettings := ah.GetLocationSettings()
		fmt.Println("b7")
		locationSettings.RLock()
		fmt.Println("b8")
		defer locationSettings.RUnlock()
		fmt.Println("b9")
		err := writeJson(locationSettings, w)
		fmt.Println("b10")
		if err != nil {
			return 500, err
		}
		fmt.Println("b11")

		return 200, nil
	}

	return 501, nil
}
