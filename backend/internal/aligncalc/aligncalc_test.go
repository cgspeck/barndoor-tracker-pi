package aligncalc

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func TestAlignCalcIgnored(t *testing.T) {
	align := models.AlignStatus{}
	locationSettings := models.LocationSettings{IgnoreAz: true}

	CalculateAlignment(&align, &locationSettings)

	e := true
	a := align.AltAligned
	if a != e {
		t.Errorf("unexpected status: got: %v, want: %v", a, e)
	}

	a = align.AzAligned
	if a != e {
		t.Errorf("unexpected status: got: %v, want: %v", a, e)
	}
}

func TestAlignCalcOnTarget(t *testing.T) {
	align := models.AlignStatus{}
	locationSettings := models.LocationSettings{IgnoreAz: true}

	CalculateAlignment(&align, &locationSettings)

	e := true
	a := align.AltAligned
	if a != e {
		t.Errorf("unexpected status: got: %v, want: %v", a, e)
	}

	a = align.AzAligned
	if a != e {
		t.Errorf("unexpected status: got: %v, want: %v", a, e)
	}
}
