package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/wireless"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type AppHandler struct {
	*models.AppContext
	H func(*models.AppContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
}

func IndexHandler(a *models.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	io.WriteString(w, fmt.Sprintf("%v", a.PreviousTime))
	return 200, nil
}

func DebugHandler(a *models.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return 500, err
	}
	io.WriteString(w, string(b))

	return 200, nil
}

func main() {
	log.Println("Barndoor Tracker Startup")

	var cmdFlags = models.CmdFlags{}

	flag.BoolVar(&cmdFlags.DisableAP, "disable-ap", false, "Disables Access Point mode")
	flag.Parse()

	previousTime := time.Now()

	context, err := CreateAppContext(previousTime, cmdFlags)
	if err != nil {
		log.Fatalf("Unable to create application context!")
	}

	err = wireless.ApplyDesiredConfiguration(context.NetworkSettings)
	if err != nil {
		log.Fatalf("Unable to apply desired network settings!\n\n%+v\n", err, context.NetworkSettings)
	}

	http.Handle("/", AppHandler{context, IndexHandler})
	http.Handle("/debug", AppHandler{context, DebugHandler})

	port := 5000
	if context.Flags.RunningAsRoot {
		port = 80
	}
	log.Printf("Start server on port %v", port)
	go http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

	for true {
		currentTime := time.Now()
		diff := currentTime.Sub(previousTime)

		if diff.Seconds() >= 10.00 {
			previousTime = currentTime
			log.Println(previousTime)
		}
	}
}
