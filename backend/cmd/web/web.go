package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/aligncalc"
	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/handlers"
	"github.com/cgspeck/barndoor-tracker-pi/internal/lsm9ds1"
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
	bus := lsm9ds1.NewMutexI2cBus(1)
	defer bus.Close()

	l, err := lsm9ds1.New(&bus)
	if err != nil {
		fmt.Printf("Error instantiating driver: %v", err)
		os.Exit(1)
	}

	if context.Arch == "arm" {
		log.Println("Begin calibration")
		l.Calibrate(true)
		log.Println("End calibration")
		log.Println("Begin Magneto calibration")
		// the next two lines can be called repeatedly until calibration looks good
		l.CalibrateMag()
		log.Printf("Mag range: %v\n", l.MagRange())
		l.LoadMagBias()
		log.Println("End Magneto calibration")
	} else {
		log.Println("Will ignore alignment")
		context.LocationSettings.IgnoreAz = true
		context.LocationSettings.IgnoreAlt = true
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case currentTime := <-ticker.C:
				diff := currentTime.Sub(previousTime)

				if diff.Milliseconds() >= 200.0 {
					if l.AccelAvailable() || l.MagAvailable(lsm9ds1.ALL_AXIS) {
						l.ReadAccel()
						l.ReadMag()
						mx, my, mz := l.M.GetReading()
						magVal := []int16{mx, my, mz}
						ax, ay, az := l.A.GetReading()
						accelVal := []int16{ax, ay, az}
						context.Lock()
						aligncalc.CalculateAlignment(
							context.AlignStatus,
							context.LocationSettings,
							accelVal,
							magVal,
						)
						context.Unlock()
					}
				}

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
