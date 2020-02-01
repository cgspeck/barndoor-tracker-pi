package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type AppHandler struct {
	*AppContext
	H func(*AppContext, http.ResponseWriter, *http.Request) (int, error)
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

func IndexHandler(a *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	io.WriteString(w, fmt.Sprintf("%v", a.PreviousTime))
	return 200, nil
}

func DebugHandler(a *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return 500, err
	}
	io.WriteString(w, string(b))

	return 200, nil
}

func main() {
	log.Println("Barndoor Tracker Startup")
	previousTime := time.Now()

	context, err := CreateAppContext(previousTime)
	if err != nil {
		log.Fatalf("Unable to create application context!")
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
