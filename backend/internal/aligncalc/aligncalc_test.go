package aligncalc

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func TestAlignCalc(t *testing.T) {

	type TestCase struct {
		ignoreAz           bool
		ignoreAlt          bool
		expectedAzAligned  bool
		expectedAltAligned bool
	}

	testCases := []TestCase{
		TestCase{
			ignoreAlt:          true,
			ignoreAz:           true,
			expectedAltAligned: true,
			expectedAzAligned:  true,
		},
		TestCase{
			ignoreAlt:          false,
			ignoreAz:           true,
			expectedAltAligned: false,
			expectedAzAligned:  true,
		},
		TestCase{
			ignoreAlt:          true,
			ignoreAz:           false,
			expectedAltAligned: true,
			expectedAzAligned:  false,
		},
	}

	for _, tt := range testCases {
		align := models.AlignStatus{}
		locationSettings := models.LocationSettings{
			IgnoreAlt: tt.ignoreAlt,
			IgnoreAz:  tt.ignoreAz,
		}

		CalculateAlignment(&align, &locationSettings)

		e := tt.expectedAltAligned
		a := align.AltAligned
		if a != e {
			t.Errorf("unexpected status: got: %v, want: %v", a, e)
		}

		e = tt.expectedAzAligned
		a = align.AzAligned
		if a != e {
			t.Errorf("unexpected status: got: %v, want: %v", a, e)
		}

	}
}
