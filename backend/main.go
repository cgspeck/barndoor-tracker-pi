package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FlagStruct struct {
	NeedsAPSettings       bool
	NeedsLocationSettings bool
}

type appContext struct {
	Flags        *FlagStruct
	PreviousTime *time.Time
}

type appHandler struct {
	*appContext
	H func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.H(ah.appContext, w, r)

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

func IndexHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	io.WriteString(w, fmt.Sprintf("%v", a.PreviousTime))
	return 200, nil
}

func main() {
	log.Println("hello world")
	previousTime := time.Now()

	var flags = FlagStruct{
		NeedsAPSettings:       false,
		NeedsLocationSettings: false,
	}

	context := &appContext{
		Flags:        &flags,
		PreviousTime: &previousTime,
	}

	http.Handle("/", appHandler{context, IndexHandler})
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
