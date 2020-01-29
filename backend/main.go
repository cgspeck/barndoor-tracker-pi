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
	PreviousTime          time.Time
}

func main() {
	log.Println("hello world")
	previousTime := time.Now()

	updateChan := make(chan time.Time)

	var flags = FlagStruct{
		NeedsAPSettings:       false,
		NeedsLocationSettings: false,
		PreviousTime:          previousTime,
	}

	go func
	http.HandleFunc("/flags", flags.flagHandler)

	// log.Fatal(http.ListenAndServe(":8080", nil))
	// go http.ListenAndServe(":8080", nil)

	for true {
		currentTime := time.Now()
		diff := currentTime.Sub(previousTime)

		if diff.Seconds() >= 2.00 {
			previousTime = currentTime
			log.Println(previousTime)
			flags.PreviousTime = previousTime
		}
	}
}

func (f FlagStruct) flagHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, fmt.Sprintf("%v", f))
}

func (f FlagStruct) updateFlags(NeedsAPSettings bool, NeedsLocationSettings bool, PreviousTime time.Time) {
	f.NeedsAPSettings = NeedsAPSettings
	f.NeedsLocationSettings = NeedsLocationSettings
	f.PreviousTime = PreviousTime
}

func (f FlagStruct) updateTime(time.Time) {
	
}