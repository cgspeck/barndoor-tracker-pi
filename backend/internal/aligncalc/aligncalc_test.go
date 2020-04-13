package aligncalc

import (
	"testing"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func TestAlignCalcIgnoreModes(t *testing.T) {

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

		CalculateAlignment(&align, &locationSettings, []float32{0, 0, 0})

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
	type TestCase struct {
		ignoreAlt          bool
		altErrorSetting    float64
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

		CalculateAlignment(&align, &locationSettings, tt.accelVal)

		e := tt.expectedAltAligned
		a := align.AltAligned
		if a != e {
			t.Errorf("unexpected status: got: %v, want: %v", a, e)
		}
	}
}
