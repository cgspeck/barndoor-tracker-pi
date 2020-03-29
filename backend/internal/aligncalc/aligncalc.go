package aligncalc

import "github.com/cgspeck/barndoor-tracker-pi/internal/models"

// CalculateAlignment returns true or false to indicate if the unit is polar aligned
func CalculateAlignment(a *models.AlignStatus, l *models.LocationSettings) {
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
}
