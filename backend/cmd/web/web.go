package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/graphlogger"
	"github.com/cgspeck/barndoor-tracker-pi/internal/runners"

	"github.com/cgspeck/barndoor-tracker-pi/internal/aligncalc"
	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/handlers"
	"github.com/cgspeck/barndoor-tracker-pi/internal/lsm9ds1"
	"github.com/cgspeck/barndoor-tracker-pi/internal/mutexi2cbus"
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

	intervalometerRunner, err := runners.NewIntervalometerRunner(5, 6, context.IntervalPeriods)

	if err != nil {
		log.Fatalf("Unable to create Intervalometer Runner: %v\n", err)
	}

	graphLogger, err := graphlogger.NewGraphLogger()
	defer graphLogger.Close()

	dewcontrollerRunner, err := runners.NewDewControllerRunner(
		context.DewControllerSettings.P,
		context.DewControllerSettings.I,
		context.DewControllerSettings.D,
		context.DewControllerSettings.TargetTemperature,
		context.DewControllerSettings.Enabled,
		25,
		context.DewControllerSettings.LoggingEnabled,
		graphLogger,
	)

	if err != nil {
		log.Fatalf("Unable to create Dew Controller Runner: %v\n", err)
	}

	http.Handle("/backend/settings/network", handlers.AppHandler{AppContext: context, H: handlers.NetworkSettingsHandler})
	http.Handle("/backend/settings/network/ap", handlers.AppHandler{AppContext: context, H: handlers.APSettingsHandler})
	// low prio: http.Handle("/backend/settings/network/profiles", handlers.AppHandler{context, ...})
	// low prio: http.Handle("/backend/settings/network/avaliable", handlers.AppHandler{context, ...})
	http.Handle("/backend/settings/intervalometer", handlers.AppHandler{
		AppContext:     context,
		H:              handlers.IntervalometerSettingsHandler,
		IntervalRunner: intervalometerRunner,
	})

	http.Handle("/backend/status/dew_controller", handlers.AppHandler{
		AppContext:          context,
		H:                   handlers.DewControllerHandler,
		DewControllerRunner: dewcontrollerRunner,
	})
	http.Handle("/backend/settings/dew_controller", handlers.AppHandler{
		AppContext:          context,
		H:                   handlers.DewControllerHandler,
		DewControllerRunner: dewcontrollerRunner,
	})
	http.Handle("/backend/settings/pid", handlers.AppHandler{
		AppContext:          context,
		H:                   handlers.DewControllerHandler,
		DewControllerRunner: dewcontrollerRunner,
	})
	http.Handle("/backend/toggle/dewcontroller", handlers.AppHandler{
		AppContext:          context,
		H:                   handlers.DewControllerHandler,
		DewControllerRunner: dewcontrollerRunner,
	})
	http.Handle("/backend/toggle/dewcontroller/logging", handlers.AppHandler{
		AppContext:          context,
		H:                   handlers.DewControllerHandler,
		DewControllerRunner: dewcontrollerRunner,
	})

	http.Handle("/backend/settings/location", handlers.AppHandler{AppContext: context, H: handlers.LocationSettingsHandler})
	http.Handle("/backend/toggle/ignoreAz", handlers.AppHandler{AppContext: context, H: handlers.LocationSettingsHandler})
	http.Handle("/backend/toggle/ignoreAlt", handlers.AppHandler{AppContext: context, H: handlers.LocationSettingsHandler})

	http.Handle("/backend/status/flags", handlers.AppHandler{AppContext: context, H: handlers.Flags})
	http.Handle("/backend/status/align", handlers.AppHandler{AppContext: context, H: handlers.AlignHandler})
	http.Handle("/backend/status/track", handlers.AppHandler{AppContext: context, H: handlers.TrackHandler})
	http.Handle("/backend/status/debug", handlers.AppHandler{AppContext: context, H: handlers.DebugHandler})

	http.Handle("/backend/toggle/intervalometer", handlers.AppHandler{AppContext: context, H: handlers.TrackHandler})
	http.Handle("/backend/track", handlers.AppHandler{AppContext: context, H: handlers.TrackHandler})

	http.HandleFunc("/backend/config.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./config.json")
	})

	port := 5000
	log.Printf("Starting server on port %v", port)

	// set up periodic callbacks
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	bus := mutexi2cbus.NewMutexI2cBus(1)
	defer bus.Close()

	l, err := lsm9ds1.New(&bus)
	if err != nil {
		fmt.Printf("Error instantiating LSM9DS1 driver: %v", err)
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

	trackerRunner := runners.NewTrackerRunner(&bus)

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

				trackerRunner.Run(currentTime, context.TrackStatus)
				intervalometerRunner.Run(currentTime, context.TrackStatus)
				dewcontrollerRunner.Run(currentTime)
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
