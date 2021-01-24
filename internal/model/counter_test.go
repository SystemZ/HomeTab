package model

import "testing"

func TestPrettyTime(t *testing.T) {
	// seconds
	r := PrettyTime(0)
	rr := "0h 0m 0s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(42)
	rr = "0h 0m 42s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(59)
	rr = "0h 0m 59s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}

	// minutes
	r = PrettyTime(60)
	rr = "0h 1m 0s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(64)
	rr = "0h 1m 4s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(119)
	rr = "0h 1m 59s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(659)
	rr = "0h 10m 59s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(3599)
	rr = "0h 59m 59s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}

	// hours
	r = PrettyTime(3600)
	rr = "1h 0m 0s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(3660)
	rr = "1h 1m 0s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(3665)
	rr = "1h 1m 5s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}
	r = PrettyTime(7199)
	rr = "1h 59m 59s"
	if r != rr {
		t.Errorf("Bad formatting, got: %s, want: %s.", r, rr)
	}

}
