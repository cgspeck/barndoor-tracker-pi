package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cgspeck/barndoor-tracker-pi/internal/runners"
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
		var newSettings runners.IntervalPeriods
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

		intervalRunner := ah.GetIntervalRunner()
		intervalRunner.Lock()
		defer intervalRunner.Unlock()
	}

	if r.Method == "GET" || r.Method == "POST" {
		intervalRunner := ah.GetIntervalRunner()
		intervalRunner.RLock()
		defer intervalRunner.RUnlock()

		intervalPeriods := intervalRunner.GetIntervalPeriods()

		err := writeJson(intervalPeriods, w)
		if err != nil {
			return 500, err
		}
		return 200, nil
	}

	return 501, nil
}
