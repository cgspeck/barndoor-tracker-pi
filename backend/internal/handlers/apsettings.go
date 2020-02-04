package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type APManagementDisabledError struct{}

func (APManagementDisabledError) Error() string {
	return "AP management is disabled"
}

func APSettingsHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method == "GET" {
		err := writeJson(ah.GetAPSettings(), w)
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

			input := make(map[string]interface{})
			err = json.Unmarshal(body, &input)
			if err != nil {
				return 500, err
			}
			err = ah.SetAPSettings(input)
			if err != nil {
				return 500, err
			}
			err = writeJson(ah.GetAPSettings(), w)
			if err != nil {
				return 500, err
			}
		} else {
			return 400, APManagementDisabledError{}
		}
	}

	return 501, nil
}
