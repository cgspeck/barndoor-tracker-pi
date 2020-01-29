package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type appHandler struct {
	*AppContext
	H func(*AppContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.H(ah.AppContext, w, r)

	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func IndexHandler(a *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	io.WriteString(w, fmt.Sprintf("%v", a.PreviousTime))
	return 200, nil
}

func DebugHandler(a *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return 500, err
	}
	io.WriteString(w, string(b))
	return 200, nil
}

func main() {
	log.Println("hello world")
	previousTime := time.Now()

	// TODO: load settings from configuration, if it exists

	var flags = FlagStruct{
		NeedsNetworkSettings:  true,
		NeedsLocationSettings: true,
	}

	var location = LocationStruct{
		Latitude:       -37.74,
		MagDeclination: 11.64,
		AzError:        2.0,
		AltError:       2.0,
		XOffset:        0,
		YOffset:        0,
		ZOffset:        0,
	}

	var alignStatus = AlignStatusStruct{
		AltAligned: true,
		AzAligned:  true,
		CurrentAz:  181.2,
		CurrentAlt: -37.4,
	}

	var networkSettings = NetworkSettingsStruct{
		AccessPointMode: true,
		APSettings: &APSettingsStruct{
			SSID:    "barndoor-tracker",
			Key:     "",
			Channel: 11,
		},
		WirelessStations: []*WirelessStation{},
	}

	context := &AppContext{
		AlignStatus:           &alignStatus,
		Flags:                 &flags,
		Location:              &location,
		PreviousTime:          &previousTime,
		NetworkSettingsStruct: &networkSettings,
	}

	http.Handle("/", appHandler{context, IndexHandler})
	http.Handle("/debug", appHandler{context, DebugHandler})
	go http.ListenAndServe(":8080", nil)

	for true {
		currentTime := time.Now()
		diff := currentTime.Sub(previousTime)

		if diff.Seconds() >= 2.00 {
			previousTime = currentTime
			log.Println(previousTime)
		}
	}
}
