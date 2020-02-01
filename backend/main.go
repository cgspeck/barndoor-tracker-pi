package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/handlers"
	"github.com/cgspeck/barndoor-tracker-pi/internal/wireless"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func IndexHandler(a *models.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	io.WriteString(w, fmt.Sprintf("%v", a.PreviousTime))
	return 200, nil
}

func main() {
	log.Println("Barndoor Tracker Startup")

	var cmdFlags = models.CmdFlags{}

	flag.BoolVar(&cmdFlags.DisableAP, "disable-ap", false, "Disables Access Point mode")
	flag.Parse()

	previousTime := time.Now()

	context, err := config.CreateAppContext(previousTime, cmdFlags)
	if err != nil {
		log.Fatalf("Unable to create application context!")
	}

	err = wireless.ApplyDesiredConfiguration(context.NetworkSettings)
	if err != nil {
		log.Fatalf("Unable to apply desired network settings: %v\n\n%+v\n", err, context.NetworkSettings)
	}

	http.Handle("/", handlers.AppHandler{context, IndexHandler})
	http.Handle("/debug", handlers.AppHandler{context, handlers.DebugHandler})

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
			context.PreviousTime = &previousTime
			log.Println(previousTime)
		}
	}
}
