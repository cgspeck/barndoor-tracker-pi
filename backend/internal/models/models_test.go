package models

import "testing"

func TestLocationSettingsEqual(t *testing.T) {
	sub := LocationSettings{
		AltError:        1,
		AzError:         2,
		Latitude:        3,
		MagDeclination:  4,
		XOffset:         5,
		YOffset:         6,
		ZOffset:         7,
		IgnoreAz: false,
	}

	same := LocationSettings{
		AltError:        1,
		AzError:         2,
		Latitude:        3,
		MagDeclination:  4,
		XOffset:         5,
		YOffset:         6,
		ZOffset:         7,
		IgnoreAz: false,
	}

	if !sub.Equals(same) {
		t.Errorf("Expected %+v to equal %+v", same, sub)
	}

	different := LocationSettings{
		AltError:        2,
		AzError:         2,
		Latitude:        3,
		MagDeclination:  4,
		XOffset:         5,
		YOffset:         6,
		ZOffset:         7,
		IgnoreAz: false,
	}

	if sub.Equals(different) {
		t.Errorf("Expected %+v not to equal %+v", different, sub)
	}
}
