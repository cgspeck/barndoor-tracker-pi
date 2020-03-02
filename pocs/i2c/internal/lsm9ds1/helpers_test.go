package lsm9ds1

import "testing"

func TestMagAvailableHelper(t *testing.T) {
	type TestCase struct {
		Name     string
		Axis     Axis
		Status   byte
		Expected bool
	}

	cases := []TestCase{
		{
			Name:     "none available no overrun",
			Axis:     ALL_AXIS,
			Status:   0x00,
			Expected: false,
		},
		{
			Name:     "none available all overrun",
			Axis:     ALL_AXIS,
			Status:   0xF0,
			Expected: false,
		},
		{
			Name:     "all available no overrun",
			Axis:     ALL_AXIS,
			Status:   0x0F,
			Expected: true,
		},
		{
			Name:     "all available all overrun",
			Axis:     ALL_AXIS,
			Status:   0xFF,
			Expected: true,
		},
		{
			Name:     "x axis, all available all overrun",
			Axis:     X_AXIS,
			Status:   0xFF,
			Expected: true,
		},
		{
			Name:     "x axis, only x available all overrun",
			Axis:     X_AXIS,
			Status:   0xF1,
			Expected: true,
		},
		{
			Name:     "x axis, only y available all overrun",
			Axis:     X_AXIS,
			Status:   0xF2,
			Expected: false,
		},
		{
			Name:     "x axis, only y & z available all overrun",
			Axis:     X_AXIS,
			Status:   0xF6,
			Expected: false,
		},
		{
			Name:     "z axis, only y & z available all overrun",
			Axis:     Z_AXIS,
			Status:   0xF6,
			Expected: true,
		},
		{
			Name:     "y axis, only y available no overrun",
			Axis:     Y_AXIS,
			Status:   0x02,
			Expected: true,
		},
		{
			Name:     "y axis, only x & z available no overrun",
			Axis:     Y_AXIS,
			Status:   0x05,
			Expected: false,
		},
	}

	for _, testCase := range cases {
		a := magAvailable(testCase.Axis, testCase.Status)

		if a != testCase.Expected {
			t.Errorf("test case %q got %v want %v", testCase.Name, a, testCase.Expected)
		}
	}
}

func TestGyroAvailableHelper(t *testing.T) {
	type TestCase struct {
		Name     string
		Status   byte
		Expected bool
	}

	cases := []TestCase{
		{
			Name:     "not available",
			Status:   0x00,
			Expected: false,
		},
		{
			Name:     "available",
			Status:   0x02,
			Expected: true,
		},
		{
			Name:     "not all available but accel is",
			Status:   0x01,
			Expected: false,
		},
	}

	for _, testCase := range cases {
		a := gyroAvailable(testCase.Status)

		if a != testCase.Expected {
			t.Errorf("test case %q got %v want %v", testCase.Name, a, testCase.Expected)
		}
	}
}

func TestAccelAvailableHelper(t *testing.T) {
	type TestCase struct {
		Name     string
		Status   byte
		Expected bool
	}

	cases := []TestCase{
		{
			Name:     "not available",
			Status:   0x00,
			Expected: false,
		},
		{
			Name:     "available",
			Status:   0x01,
			Expected: true,
		},
		{
			Name:     "not all available but gyro is",
			Status:   0x02,
			Expected: false,
		},
	}

	for _, testCase := range cases {
		a := accelAvailable(testCase.Status)

		if a != testCase.Expected {
			t.Errorf("test case %q got %v want %v", testCase.Name, a, testCase.Expected)
		}
	}
}
