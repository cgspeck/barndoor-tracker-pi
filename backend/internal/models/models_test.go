package models

import "testing"

func TestLocationSettingsEqual(t *testing.T) {
	sub := LocationSettings{
		AltError:       1,
		AzError:        2,
		Latitude:       3,
		MagDeclination: 4,
		XOffset:        5,
		YOffset:        6,
		ZOffset:        7,
		IgnoreAz:       false,
	}

	same := LocationSettings{
		AltError:       1,
		AzError:        2,
		Latitude:       3,
		MagDeclination: 4,
		XOffset:        5,
		YOffset:        6,
		ZOffset:        7,
		IgnoreAz:       false,
	}

	if !sub.Equals(same) {
		t.Errorf("Expected %+v to equal %+v", same, sub)
	}

	different := LocationSettings{
		AltError:       2,
		AzError:        2,
		Latitude:       3,
		MagDeclination: 4,
		XOffset:        5,
		YOffset:        6,
		ZOffset:        7,
		IgnoreAz:       false,
	}

	if sub.Equals(different) {
		t.Errorf("Expected %+v not to equal %+v", different, sub)
	}
}

func TestTrackStatusProcessTrackCommand(t *testing.T) {
	type TestCase struct {
		previousState string
		input         string
		expectedError bool
		expectedState string
	}

	testCases := []TestCase{
		TestCase{
			previousState: "Idle",
			input:         "home",
			expectedError: false,
			expectedState: "Homing Requested",
		},
		TestCase{
			previousState: "Finished",
			input:         "home",
			expectedError: false,
			expectedState: "Homing Requested",
		},
		TestCase{
			previousState: "Homing",
			input:         "home",
			expectedError: true,
			expectedState: "Homing",
		},
		TestCase{
			previousState: "Homed",
			input:         "track",
			expectedError: false,
			expectedState: "Tracking Requested",
		},
		TestCase{
			previousState: "Idle",
			input:         "track",
			expectedError: true,
			expectedState: "Idle",
		},
		TestCase{
			previousState: "Tracking",
			input:         "stop",
			expectedError: false,
			expectedState: "Idle",
		},
		TestCase{
			previousState: "Homing",
			input:         "stop",
			expectedError: true,
			expectedState: "Homing",
		},
		TestCase{
			previousState: "Finished",
			input:         "fauxcommand",
			expectedError: true,
			expectedState: "Finished",
		},
	}

	for _, tt := range testCases {
		ts := TrackStatus{
			State:         tt.previousState,
			PreviousState: "foo",
		}

		actualState, actualError := ts.ProcessTrackCommand(tt.input)

		if ts.State != tt.expectedState {
			t.Errorf("unexpected: got: %q, want: %q, input: %q, previous: %q", ts.State, tt.expectedState, tt.input, tt.previousState)
		}

		if !tt.expectedError {
			if actualState != tt.expectedState {
				t.Errorf("unexpected: got: %q, want: %q, input: %q, previous: %q", actualState, tt.expectedState, tt.input, tt.previousState)
			}
			if ts.PreviousState != tt.previousState {
				t.Errorf("unexpected PreviousState: got: %q, want: %q, input: %q, previous: %q", ts.PreviousState, ts.PreviousState, tt.input, tt.previousState)
			}
			if actualError != nil {
				t.Errorf("unexpected error: got: %v, input: %q, previous: %q", actualError, tt.input, tt.previousState)
			}
		} else {
			if actualError == nil {
				t.Errorf("expected error but got none!: input: %q, previous: %q", tt.input, tt.previousState)
			}

			if ts.PreviousState != "foo" {
				t.Errorf("unexpected PreviousState: got: %q, want: %q, input: %q, previous: %q", ts.PreviousState, "foo", tt.input, tt.previousState)
			}
		}
	}
}

func TestTrackStatusProcessArduinoStateChange(t *testing.T) {
	type TestCase struct {
		previousState string
		input         string
		expectedError bool
		expectedState string
	}

	testCases := []TestCase{
		TestCase{
			previousState: "Idle",
			input:         "Homing",
			expectedError: false,
			expectedState: "Homing",
		},
		TestCase{
			previousState: "Homing Requested",
			input:         "Homing",
			expectedError: false,
			expectedState: "Homing",
		},
		TestCase{
			previousState: "Finished",
			input:         "Homing",
			expectedError: false,
			expectedState: "Homing",
		},
		TestCase{
			previousState: "Tracking",
			input:         "Homing",
			expectedError: true,
			expectedState: "Tracking",
		},
		TestCase{
			previousState: "Homed",
			input:         "Tracking",
			expectedError: false,
			expectedState: "Tracking",
		},
		TestCase{
			previousState: "Idle",
			input:         "Tracking",
			expectedError: true,
			expectedState: "Idle",
		},
		TestCase{
			previousState: "Tracking",
			input:         "Finished",
			expectedError: false,
			expectedState: "Finished",
		},
		TestCase{
			previousState: "Idle",
			input:         "Finished",
			expectedError: true,
			expectedState: "Idle",
		},
		TestCase{
			previousState: "Finished",
			input:         "FauxCommand",
			expectedError: true,
			expectedState: "Finished",
		},
	}

	for _, tt := range testCases {
		ts := TrackStatus{
			State:         tt.previousState,
			PreviousState: "foo",
		}

		actualState, actualError := ts.ProcessArduinoStateChange(tt.input)

		if ts.State != tt.expectedState {
			t.Errorf("unexpected: got: %q, want: %q, input: %q, previous: %q", ts.State, tt.expectedState, tt.input, tt.previousState)
		}

		if !tt.expectedError {
			if actualState != tt.expectedState {
				t.Errorf("unexpected: got: %q, want: %q, input: %q, previous: %q", actualState, tt.expectedState, tt.input, tt.previousState)
			}
			if ts.PreviousState != tt.previousState {
				t.Errorf("unexpected PreviousState: got: %q, want: %q, input: %q, previous: %q", ts.PreviousState, ts.PreviousState, tt.input, tt.previousState)
			}
			if actualError != nil {
				t.Errorf("unexpected error: got: %v, input: %q, previous: %q", actualError, tt.input, tt.previousState)
			}
		} else {
			if actualError == nil {
				t.Errorf("expected error but got none!: input: %q, previous: %q", tt.input, tt.previousState)
			}

			if ts.PreviousState != "foo" {
				t.Errorf("unexpected PreviousState: got: %q, want: %q, input: %q, previous: %q", ts.PreviousState, "foo", tt.input, tt.previousState)
			}
		}
	}
}
