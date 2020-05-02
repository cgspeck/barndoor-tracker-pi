package aligncalc

import (
	"math"
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func floatEquals(a, b float64) bool {
	ba := math.Float64bits(a)
	bb := math.Float64bits(b)
	diff := ba - bb
	if diff < 0 {
		diff = -diff
	}
	// accept one bit difference
	return diff < 2
}

func TestAlignCalcIgnoreModes(t *testing.T) {
	nilReading := []int16{0, 0, 0}

	type TestCase struct {
		desc               string
		ignoreAz           bool
		ignoreAlt          bool
		expectedAzAligned  bool
		expectedAltAligned bool
	}

	testCases := []TestCase{
		TestCase{
			desc:               "ignore both",
			ignoreAlt:          true,
			ignoreAz:           true,
			expectedAltAligned: true,
			expectedAzAligned:  true,
		},
		TestCase{
			desc:               "ignore Alt",
			ignoreAlt:          false,
			ignoreAz:           true,
			expectedAltAligned: false,
			expectedAzAligned:  true,
		},
		TestCase{
			desc:               "ignore AZ",
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
			Latitude:  -37.813,
		}

		CalculateAlignment(&align, &locationSettings, nilReading, nilReading)

		e := tt.expectedAltAligned
		a := align.AltAligned
		if a != e {
			t.Errorf("unexpected status: got: %v, want: %v, case: %q", a, e, tt.desc)
		}

		e = tt.expectedAzAligned
		a = align.AzAligned
		if a != e {
			t.Errorf("unexpected status: got: %v, want: %v, case: %q", a, e, tt.desc)
		}
	}
}

func TestAlignCalcAltitude(t *testing.T) {
	nilReading := []int16{0, 0, 0}

	locationSettings := models.LocationSettings{
		AltError:  1.0,
		IgnoreAlt: false,
		Latitude:  -44.35,
	}

	type TestCase struct {
		desc               string
		expectedCurrentAlt float64
		expectedAltAligned bool
		accelVal           []int16
	}

	testCases := []TestCase{
		TestCase{
			desc:               "pitched down into the ground",
			expectedAltAligned: false,
			expectedCurrentAlt: -44.739905971358276,
			accelVal:           []int16{10587, -10573, 1533},
		},
		TestCase{
			desc:               "pitched up on target",
			expectedAltAligned: true,
			expectedCurrentAlt: 44.739905971358276,
			accelVal:           []int16{-10587, -10573, 1533},
		},
		TestCase{
			desc:               "exceeds setting by 1.01, which is > error",
			expectedAltAligned: false,
			expectedCurrentAlt: 45.364838727325385,
			accelVal:           []int16{-10587, -10340, 1533},
		},
	}

	for _, tt := range testCases {
		align := models.AlignStatus{}

		CalculateAlignment(&align, &locationSettings, tt.accelVal, nilReading)

		if tt.expectedCurrentAlt != 0 {
			e := tt.expectedCurrentAlt
			a := align.CurrentAlt
			if !floatEquals(a, e) {
				t.Errorf("unexpected CurrentAlt: got: %v, want: %v, case: %q", a, e, tt.desc)
			}
		}

		e := tt.expectedAltAligned
		a := align.AltAligned
		if a != e {
			t.Errorf("unexpected AltAligned: got: %v, want: %v, case: %q", a, e, tt.desc)
		}
	}
}

func TestAlignCalcAzimuthSouthernHemisphere(t *testing.T) {
	nilReading := []int16{0, 0, 0}

	locationSettings := models.LocationSettings{
		AzError:        0.1,
		IgnoreAz:       false,
		Latitude:       -44.35,
		MagDeclination: 0,
	}

	type TestCase struct {
		desc                string
		expectedAzAligned   bool
		expectedCalcHeading float64
		magVal              []int16
	}

	testCases := []TestCase{
		TestCase{
			desc:                "on target",
			expectedAzAligned:   true,
			expectedCalcHeading: 180,
			magVal:              []int16{},
		},
	}

	for _, tt := range testCases {
		align := models.AlignStatus{}
		CalculateAlignment(&align, &locationSettings, nilReading, tt.magVal)

		eStatus := tt.expectedAzAligned
		aStatus := align.AzAligned
		if aStatus != eStatus {
			t.Errorf("unexpected status: got: %v, want: %v, case: %q", aStatus, eStatus, tt.desc)
		}

		eVal := tt.expectedCalcHeading
		aVal := align.CurrentAz
		if aVal != eVal {
			t.Errorf("unexpected calculation: got: %v, want: %v, case: %q", aVal, eVal, tt.desc)
		}
	}
}

func TestAlignCalcAzimuthNorthernHemisphere(t *testing.T) {}

func TestAlignCalcAzimuthWithDeclination(t *testing.T) {}

func TestCalculateHeading(t *testing.T) {
	testCases := []struct {
		mx          int16
		my          int16
		declination float64
		expected    float64
	}{
		{
			mx:       60,
			my:       0,
			expected: 0,
		},
		{
			mx:       47,
			my:       1,
			expected: 1.2188752351312977,
		},
		{
			mx:       -65,
			my:       1,
			expected: 179.11859600341788,
		},
		{
			mx:       -60,
			my:       0,
			expected: 180,
		},
		{
			mx:          -60,
			my:          0,
			declination: 10,
			expected:    190,
		},
		{
			mx:          -60,
			my:          0,
			declination: -10,
			expected:    170,
		},
		{
			mx:          60,
			my:          0,
			declination: -10,
			expected:    350,
		},
	}

	for _, tt := range testCases {
		a := calculateHeading([]int16{tt.mx, tt.my, 0}, tt.declination)
		if !floatEquals(tt.expected, a) {
			t.Errorf("unexpected result: got: %v want: %v x: %v y: %v",
				a, tt.expected, tt.mx, tt.my,
			)
		}
	}
}
