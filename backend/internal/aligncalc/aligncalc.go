package aligncalc

import (
	"fmt"
	"math"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

// CalculateAlignment returns true or false to indicate if the unit is polar aligned
func CalculateAlignment(a *models.AlignStatus, l *models.LocationSettings, accelVal []int16, magVal []int16) {
	l.RLock()
	IgnoreAz := l.IgnoreAz
	IgnoreAlt := l.IgnoreAlt
	l.RUnlock()

	a.Lock()
	defer a.Unlock()

	if IgnoreAz {
		a.AzAligned = true
	}

	if IgnoreAlt {
		a.AltAligned = true
	}

	if IgnoreAlt && IgnoreAz {
		return
	}

	if !IgnoreAlt {
		fmt.Println("accelVal")
		fmt.Println(accelVal)
		accelVal64 := []float64{float64(accelVal[0]), float64(accelVal[1]), float64(accelVal[2])}
		altVal := math.Atan2(-accelVal64[0], math.Sqrt(accelVal64[1]*accelVal64[1]+accelVal64[2]*accelVal64[2]))
		fmt.Println("accelVal64")
		fmt.Println(accelVal64)
		fmt.Println("radians")
		fmt.Println(altVal)
		altVal = altVal * 180 / math.Pi
		fmt.Println("altVal")
		fmt.Println(altVal)
		a.CurrentAlt = altVal

		absLat := math.Abs(l.Latitude)
		minAlt := absLat - l.AltError
		maxAlt := absLat + l.AltError
		a.AltAligned = (altVal >= minAlt && altVal <= maxAlt)
	}
}
