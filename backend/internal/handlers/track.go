package handlers

import (
	"net/http"
)

// type APManagementDisabledError struct{}

// func (APManagementDisabledError) Error() string {
// 	return "AP management is disabled"
// }

func TrackHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	// if r.Method == "POST" {
	// 	networkSettings := ah.GetNetworkSettings()
	// 	networkSettings.RLock()
	// 	managementEnabled := networkSettings.ManagementEnabled
	// 	networkSettings.RUnlock()

	// 	if managementEnabled {
	// 		body, err := ioutil.ReadAll(r.Body)
	// 		if err != nil {
	// 			return 500, err
	// 		}

	// 		input := make(map[string]interface{})
	// 		err = json.Unmarshal(body, &input)
	// 		if err != nil {
	// 			return 500, err
	// 		}

	// 		apSettings := ah.GetAPSettings()
	// 		apSettings.Lock()
	// 		err = ah.SetAPSettings(input)
	// 		apSettings.Unlock()

	// 		if err != nil {
	// 			return 500, err
	// 		}
	// 	} else {
	// 		return 400, APManagementDisabledError{}
	// 	}
	// }

	if r.Method == "GET" || r.Method == "POST" {
		trackStatus := ah.GetTrackStatus()
		trackStatus.RLock()
		defer trackStatus.RUnlock()
		err := writeJson(trackStatus, w)
		if err != nil {
			return 500, err
		}
		return 200, nil
	}

	return 501, nil
}
