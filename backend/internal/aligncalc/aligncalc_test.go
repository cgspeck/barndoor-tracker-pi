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

		CalculateAlignment(&align, &locationSettings, nilReading, nilReading)

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

func TestAlignCalcAltitude(t *testing.T) {
	nilReading := []int16{0, 0, 0}

	type TestCase struct {
		ignoreAlt          bool
		latitudeSetting    float64
		calcLatitude       float64
		expectedAltAligned bool
		accelVal           []int16
		AltError           float64
	}

	testCases := []TestCase{
		TestCase{
			ignoreAlt:          false,
			expectedAltAligned: false,
			AltError:           2.0,
			latitudeSetting:    -44.35,
			calcLatitude:       -44.739905971358276,
			accelVal:           []int16{10587, -10573, 1533},
		},
		// TestCase{
		// 	ignoreAlt:          true,
		// 	expectedAltAligned: true,
		// 	AltError:           0.1,
		// 	accelVal:           []int16{0, 0, 0},
		// },
		// TestCase{
		// 	ignoreAlt:          false,
		// 	expectedAltAligned: false,
		// 	AltError:           0.1,
		// 	accelVal:           []int16{0, 0, 0},
		// },
		// TestCase{
		// 	ignoreAlt:          false,
		// 	expectedAltAligned: true,
		// 	AltError:           0.1,
		// 	latitudeSetting:    42.34,
		// 	calcLatitude:       44.61,
		// 	accelVal:           []int16{-10536, -10566, 1541},
		// },
		// TestCase{
		// 	ignoreAlt:          false,
		// 	expectedAltAligned: true,
		// 	AltError:           0.1,
		// 	latitudeSetting:    -44.61,
		// 	calcLatitude:       44.61,
		// 	accelVal:           []int16{-10536, -10566, 1541},
		// },
		// TestCase{
		// 	ignoreAlt:          false,
		// 	expectedAltAligned: true,
		// 	AltError:           2.0,
		// 	latitudeSetting:    44.34,
		// 	calcLatitude:       44.61,
		// 	accelVal:           []int16{-10536, -10566, 1541},
		// },
		// TestCase{
		// 	ignoreAlt:          false,
		// 	expectedAltAligned: true,
		// 	AltError:           2.0,
		// 	latitudeSetting:    -44.34,
		// 	calcLatitude:       44.61,
		// 	accelVal:           []int16{-10536, -10566, 1541},
		// },
		// TestCase{
		// 	ignoreAlt:          false,
		// 	expectedAltAligned: false,
		// 	AltError:           2.0,
		// 	latitudeSetting:    -44.35,
		// 	calcLatitude:       44.61,
		// 	accelVal:           []int16{-10536, -10566, 1541},
		// },
	}

	for i, tt := range testCases {
		align := models.AlignStatus{}
		locationSettings := models.LocationSettings{
			AltError:  tt.AltError,
			IgnoreAlt: tt.ignoreAlt,
			Latitude:  tt.latitudeSetting,
		}

		CalculateAlignment(&align, &locationSettings, tt.accelVal, nilReading)

		if tt.calcLatitude != 0 {
			e := tt.calcLatitude
			a := align.CurrentAlt
			if !floatEquals(a, e) {
				t.Errorf("case %v unexpected CurrentAlt value: got: %v, want: %v", i+1, a, e)
			}
		}

		e := tt.expectedAltAligned
		a := align.AltAligned
		if a != e {
			t.Errorf("case %v unexpected AltAligned status: got: %v, want: %v", i+1, a, e)
		}
	}
}

func TestAlignCalcAzimuth(t *testing.T) {
	nilReading := []int16{0, 0, 0}

	type TestCase struct {
		azErrorSetting      float64
		declinationSetting  float64
		expectedAzAligned   bool
		expectedCalcHeading float64
		ignoreAz            bool
		latitudeSetting     float64
		magVal              []int16
	}

	testCases := []TestCase{
		TestCase{
			ignoreAz:          true,
			expectedAzAligned: true,
		},
		TestCase{
			azErrorSetting:      0.1,
			declinationSetting:  0.0,
			expectedAzAligned:   true,
			expectedCalcHeading: 10,
			ignoreAz:            false,
			latitudeSetting:     -37,
			magVal:              []int16{},
		},
	}

	for _, tt := range testCases {
		align := models.AlignStatus{}
		locationSettings := models.LocationSettings{
			AzError:  tt.azErrorSetting,
			IgnoreAz: tt.ignoreAz,
			Latitude: tt.latitudeSetting,
		}

		CalculateAlignment(&align, &locationSettings, nilReading, tt.magVal)

		eStatus := tt.expectedAzAligned
		aStatus := align.AzAligned
		if aStatus != eStatus {
			t.Errorf("unexpected status: got: %v, want: %v", aStatus, eStatus)
		}

		eVal := tt.expectedCalcHeading
		aVal := align.CurrentAz
		if aVal != eVal {
			t.Errorf("unexpected calculation: got: %v, want: %v", aVal, eVal)
		}
	}
}
