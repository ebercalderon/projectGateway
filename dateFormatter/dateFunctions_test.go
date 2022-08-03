package dateFormatter

import "testing"

func TestGetStartOfDay(t *testing.T) {
	var got int64 = GetStartOfDay(1656363850000)
	var want int64 = 1656288000000

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestGetEndOfDay(t *testing.T) {
	var got int64 = GetEndOfDay(1656363850000)
	var want int64 = 1656374399000

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
