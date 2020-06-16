package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UnrecognisedPathError struct {
	Path string
}

func (u UnrecognisedPathError) Error() string {
	return fmt.Sprintf("Unrecognised path: %q", u.Path)
}

type BadRequestError struct{}

func (_ BadRequestError) Error() string {
	return "Bad Request"
}

type CouldNotCastToBoolError struct {
	val interface{}
}

func (e CouldNotCastToBoolError) Error() string {
	return fmt.Sprintf("Could not cast %q to bool", e.val)
}

func handleTrackCommand(ah IAppHandler, command string) error {
	trackStatus := ah.GetTrackStatus()
	trackStatus.Lock()
	defer trackStatus.Unlock()

	_, err := trackStatus.ProcessTrackCommand(command)
	return err
}

func TrackHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	/*
		POST handles these routes:

		"/backend/toggle/intervalometer", { "enabled": bool }
		"/backend/toggle/dewcontroller", { "enabled": bool }
		"/backend/track" { "command": "val"} <- need to check this against current state!

	*/
	if r.Method == "POST" {
		path := r.URL.Path

		if !(path == "/backend/track" || path == "/backend/toggle/intervalometer" || path == "/backend/toggle/dewcontroller") {
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
		case "/backend/track":
			command, ok := input["command"]
			if ok {
				err = handleTrackCommand(ah, fmt.Sprintf("%v", command))
			} else {
				err = BadRequestError{}
			}
		case "/backend/toggle/intervalometer", "/backend/toggle/dewcontroller":
			iEnabled, ok := input["enabled"]

			if !ok {
				err = BadRequestError{}
			}

			bEnabled, ok := iEnabled.(bool)

			if !ok {
				err = CouldNotCastToBoolError{iEnabled}
			}

			trackStatus := ah.GetTrackStatus()

			if path == "/backend/toggle/intervalometer" {
				trackStatus.IntervolmeterEnabled = bEnabled
			} else {
				trackStatus.DewControllerEnabled = bEnabled
			}
		}
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
