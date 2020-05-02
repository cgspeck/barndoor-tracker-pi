package aligncalc

import (
	"math"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func calculateHeading(magVal []int16, declination float64) float64 {
	magVal64 := []float64{float64(magVal[0]), float64(magVal[1]), float64(magVal[2])}
	mx := magVal64[0]
	my := magVal64[1]

	// https://github.com/sparkfun/SparkFun_LSM9DS1_Arduino_Library/blob/master/examples/LSM9DS1_Basic_SPI/LSM9DS1_Basic_SPI.ino#L226
	// float heading;
	// if (my == 0)
	//     heading = (mx < 0) ? PI : 0;
	//   else
	//     heading = atan2(mx, my);

	//   heading -= DECLINATION * PI / 180;

	//   if (heading > PI) heading -= (2 * PI);
	//   else if (heading < -PI) heading += (2 * PI);

	//   // Convert everything from radians to degrees:
	//   heading *= 180.0 / PI;
	var heading float64 = 0

	if my == 0 {
		if mx < 0 {
			heading = math.Pi
		}
	} else {
		heading = math.Atan2(my, mx)
	}

	heading -= declination * math.Pi / 180

	if heading > math.Pi {
		heading -= (2 * math.Pi)
	} else if heading < -math.Pi {
		heading += (2 * math.Pi)
	}

	heading *= 180.0 / math.Pi

	return heading
}

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
		accelVal64 := []float64{float64(accelVal[0]), float64(accelVal[1]), float64(accelVal[2])}
		altVal := math.Atan2(-accelVal64[0], math.Sqrt(accelVal64[1]*accelVal64[1]+accelVal64[2]*accelVal64[2]))
		altVal = altVal * 180 / math.Pi
		a.CurrentAlt = altVal

		absLat := math.Abs(l.Latitude)
		minAlt := absLat - l.AltError
		maxAlt := absLat + l.AltError
		a.AltAligned = (altVal >= minAlt && altVal <= maxAlt)
	}
}
