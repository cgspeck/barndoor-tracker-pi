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

	http.Handle("/settings/network", handlers.AppHandler{AppContext: context, H: handlers.NetworkSettingsHandler})
	http.Handle("/settings/network/ap", handlers.AppHandler{AppContext: context, H: handlers.APSettingsHandler})
	// low prio: http.Handle("/settings/network/profiles", handlers.AppHandler{context, ...})
	// low prio: http.Handle("/settings/network/avaliable", handlers.AppHandler{context, ...})

	http.Handle("/settings/location", handlers.AppHandler{AppContext: context, H: handlers.LocationSettingsHandler})

	http.Handle("/status/flags", handlers.AppHandler{AppContext: context, H: handlers.Flags})
	http.Handle("/status/align", handlers.AppHandler{AppContext: context, H: handlers.AlignHandler})
	// http.Handle("/status/track", handlers.AppHandler{context, ...})
	http.Handle("/status/debug", handlers.AppHandler{AppContext: context, H: handlers.DebugHandler})

	// location of the React/Preact Frontend
	static := "../frontend/build"
	if context.Arch == "arm" {
		static = "html"
	}
	log.Printf("Serving static content from %v", static)
	fs := http.FileServer(http.Dir(static))
	http.Handle("/", fs)

	http.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./config.json")
	})

	port := 5000
	if context.Flags.RunningAsRoot {
		port = 80
	}
	log.Printf("Starting server on port %v", port)

	// set up periodic callbacks
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case currentTime := <-ticker.C:
				diff := currentTime.Sub(previousTime)

				if diff.Seconds() >= 10.00 {
					previousTime = currentTime
					context.Lock()
					context.Time = &previousTime
					context.Unlock()
					log.Println(previousTime)
				}
			}
		}
	}()

	defer sendDone(done)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func sendDone(ch chan<- bool) {
	ch <- true
}
