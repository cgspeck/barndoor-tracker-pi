package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type InvalidBulbInterval struct {
	Value int
}

func (e InvalidBulbInterval) Error() string {
	return fmt.Sprintf("Invalid BulbInterval: %d", e.Value)
}

type RestBulbInterval struct {
	Value int
}

func (e RestBulbInterval) Error() string {
	return fmt.Sprintf("Invalid RestInterval: %d", e.Value)
}

func IntervalometerSettingsHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	/*
		GET/POST handles these routes:

		"/backend/settings/intervalometer", {"bulbInterval": 60, "restInterval": 61}

	*/
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading stream: %v\n", err)
			return 500, err
		}
		var newSettings models.IntervalPeriods
		err = json.Unmarshal(body, &newSettings)
		if err != nil {
			log.Printf("Error unmarshalling json: %v\n", err)
			return 500, err
		}

		if newSettings.BulbTimeSeconds < 1 {
			return 400, InvalidBulbInterval{newSettings.BulbTimeSeconds}
		}

		if newSettings.RestTimeSeconds < 1 {
			return 400, RestBulbInterval{newSettings.RestTimeSeconds}
		}

		err = ah.SaveIntervalPeriods(&newSettings)
		if err != nil {
			log.Printf("Error applying update: %v\n", err)
			return 500, err
		}
	}

	if r.Method == "GET" || r.Method == "POST" {
		intervalPeriods := ah.GetIntervalPeriods()

		err := writeJson(intervalPeriods, w)
		if err != nil {
			return 500, err
		}
		return 200, nil
	}

	return 501, nil
}
