package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func extractFloat64(hashMap map[string]interface{}, key string) (float64, error) {
	v, ok := hashMap[key]

	if !ok {
		return 0, KeyNotFoundError{hashMap: hashMap, k: key}
	}

	f, ok := v.(float64)

	if !ok {
		return 0, CouldNotCastToFloat64Error{v}
	}

	return f, nil
}

func extractPID(hashMap map[string]interface{}) (float64, float64, float64, error) {
	res := map[string]float64{
		"p": 0,
		"i": 0,
		"d": 0,
	}

	var err error

	for k := range res {
		res[k], err = extractFloat64(hashMap, k)
		if err != nil {
			break
		}
	}

	if err != nil {
		return 0, 0, 0, err
	}

	return res["p"], res["i"], res["d"], nil
}

func DewControllerHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	/*
		POST handles these routes:
		/backend/settings/dew_controller (set targetTemp)
		/backend/settings/pid
		/backend/toggle/dewcontroller

	*/

	if r.Method == "POST" {
		path := r.URL.Path

		if !(path == "/backend/settings/dew_controller" || path == "/backend/settings/pid" || path == "/backend/toggle/dewcontroller") {
			return 404, UnrecognisedPathError{path}
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return 500, err
		}
		input := make(map[string]interface{})
		err = json.Unmarshal(body, &input)
		if err != nil {
			return 500, err
		}

		switch path {
		case "/backend/settings/dew_controller":
			v, err := extractFloat64(input, "targetTemp")

			if err != nil {
				return 500, err
			}

			err = ah.SetTargetTemperature(v)

			if err != nil {
				return 500, err
			}

		case "/backend/settings/pid":
			p, i, d, err := extractPID(input)
			if err != nil {
				return 500, err
			}

			ah.SetPID(p, i, d)

		case "/backend/toggle/dewcontroller":
			iEnabled, ok := input["enabled"]

			if !ok {
				err = BadRequestError{}
			}

			bEnabled, ok := iEnabled.(bool)

			if !ok {
				err = CouldNotCastToBoolError{iEnabled}
			}

			err = ah.SetDewControllerEnabled(bEnabled)

			if err != nil {
				return 500, err
			}

			fmt.Fprintf(w, "{ \"enabled\" : %v }", iEnabled)
			return 200, nil
		}
	}

	if r.Method == "GET" || r.Method == "POST" {
		dewControllerStatus := ah.GetDewControllerStatus()

		err := writeJson(dewControllerStatus, w)
		if err != nil {
			return 500, err
		}
		return 200, nil
	}

	return 501, nil
}
