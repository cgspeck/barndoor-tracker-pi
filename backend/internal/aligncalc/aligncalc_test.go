package aligncalc

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func TestAlignCalcIgnoreModes(t *testing.T) {
	nilReading := []float32{0, 0, 0}

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
	nilReading := []float32{0, 0, 0}

	type TestCase struct {
		ignoreAlt          bool
		latitudeSetting    float64
		calcLatitude       float64
		expectedAltAligned bool
		accelVal           []float32
		AltError           float64
	}

	testCases := []TestCase{
		TestCase{
			ignoreAlt:          true,
			expectedAltAligned: true,
			AltError:           0.1,
			accelVal:           []float32{0, 0, 0},
		},
		TestCase{
			ignoreAlt:          false,
			expectedAltAligned: false,
			AltError:           0.1,
			accelVal:           []float32{0, 0, 0},
		},
		TestCase{
			ignoreAlt:          false,
			expectedAltAligned: true,
			AltError:           0.1,
			latitudeSetting:    42.34,
			calcLatitude:       42.34,
			accelVal:           []float32{-0.771406, 0.511119, 0.48263198},
		},
		TestCase{
			ignoreAlt:          false,
			expectedAltAligned: true,
			AltError:           0.1,
			latitudeSetting:    -42.34,
			calcLatitude:       42.34,
			accelVal:           []float32{-0.771406, 0.511119, 0.48263198},
		},
		TestCase{
			ignoreAlt:          false,
			expectedAltAligned: true,
			AltError:           2.0,
			latitudeSetting:    44.34,
			calcLatitude:       42.34,
			accelVal:           []float32{-0.771406, 0.511119, 0.48263198},
		},
		TestCase{
			ignoreAlt:          false,
			expectedAltAligned: true,
			AltError:           2.0,
			latitudeSetting:    -44.34,
			calcLatitude:       42.34,
			accelVal:           []float32{-0.771406, 0.511119, 0.48263198},
		},
		TestCase{
			ignoreAlt:          false,
			expectedAltAligned: false,
			AltError:           2.0,
			latitudeSetting:    -44.35,
			calcLatitude:       42.34,
			accelVal:           []float32{-0.771406, 0.511119, 0.48263198},
		},
	}

	for _, tt := range testCases {
		align := models.AlignStatus{}
		locationSettings := models.LocationSettings{
			AltError:  tt.AltError,
			IgnoreAlt: tt.ignoreAlt,
			Latitude:  tt.latitudeSetting,
		}

		CalculateAlignment(&align, &locationSettings, tt.accelVal, nilReading)

		e := tt.expectedAltAligned
		a := align.AltAligned
		if a != e {
			t.Errorf("unexpected status: got: %v, want: %v", a, e)
		}
	}
}

func TestAlignCalcAzimuth(t *testing.T) {
	nilReading := []float32{0, 0, 0}

	type TestCase struct {
		azErrorSetting      float64
		declinationSetting  float64
		expectedAzAligned   bool
		expectedCalcHeading float64
		ignoreAz            bool
		latitudeSetting     float64
		magVal              []float32
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
			magVal:              []float32{},
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
