package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// type APManagementDisabledError struct{}

// func (APManagementDisabledError) Error() string {
// 	return "AP management is disabled"
// }

type UnrecognisedPathError struct {
	Path string
}

func (u UnrecognisedPathError) Error() string {
	return fmt.Sprintf("Unrecognised path: %q", u.Path)
}

func TrackHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	/*
		POST handles these routes:

		"/toggle/intervalometer", { "enabled": bool }
		"/toggle/dewcontroller", { "enabled": bool }
		"/track" { "command": "val"} <- need to check this against current state!

	*/
	if r.Method == "POST" {
		path := r.URL.Path

		if !(path == "/track" || path == "/toggle/intervalometer" || path == "/toggle/dewcontroller") {
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
		case "/track":
			err = handleTrackCommand()
		}

		// newState, err

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
	}

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
