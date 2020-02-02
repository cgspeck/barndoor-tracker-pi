package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/handlers"
	"github.com/cgspeck/barndoor-tracker-pi/internal/wireless"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

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

	http.Handle("/", handlers.AppHandler{AppContext: context, H: handlers.IndexHandler})

	http.Handle("/settings/network", handlers.AppHandler{AppContext: context, H: handlers.NetworkSettingsHandler})
	// http.Handle("/settings/network/ap", handlers.AppHandler{context, ...})
	// low prio: http.Handle("/settings/network/profiles", handlers.AppHandler{context, ...})
	// low prio: http.Handle("/settings/network/avaliable", handlers.AppHandler{context, ...})

	// http.Handle("/settings/location", handlers.AppHandler{context, ...})

	// http.Handle("/status/flags", handlers.AppHandler{context, ...})
	// http.Handle("/status/align", handlers.AppHandler{context, ...})
	http.Handle("/status/debug", handlers.AppHandler{AppContext: context, H: handlers.DebugHandler})

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
			context.Time = &previousTime
			log.Println(previousTime)
		}
	}
}
