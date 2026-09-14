package util

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	cases := []struct {
		value string
		start int
		end   int
	}{
		{"2026-09-13", 20260913, 20260913},
		{"13-09-2026", 20260913, 20260913},
		{"2026/09/13", 20260913, 20260913},
		{"13/09/2026", 20260913, 20260913},
		{"2026-09", 20260900, 20260999},
		{"09-2026", 20260900, 20260999},
		{"2026/09", 20260900, 20260999},
		{"09/2026", 20260900, 20260999},
		{"2026", 20260000, 20269999},
	}

	for _, c := range cases {
		start, end, err := ParseDate(c.value)
		if err != nil {
			t.Fatalf("%q: unexpected error %v", c.value, err)
		}
		if start != c.start || end != c.end {
			t.Fatalf("%q: expected [%d %d], got [%d %d]", c.value, c.start, c.end, start, end)
		}
	}
}

func TestParseDateYearPositionDecidesOrder(t *testing.T) {
	forms := []string{"2026-09-13", "13-09-2026", "2026/09/13", "13/09/2026"}

	for _, form := range forms {
		start, _, err := ParseDate(form)
		if err != nil {
			t.Fatalf("%q: unexpected error %v", form, err)
		}

		if start != 20260913 {
			t.Fatalf("%q: expected 20260913, got %d", form, start)
		}
	}
}

func TestParseDateRejectsInvalid(t *testing.T) {
	for _, value := range []string{
		"", "yesterday", "2026-13-01", "2026-00-01", "2026-09-32",
		"2026-09-13-01", "abc-09-13", "0000", "99999",
		"13-09-26", "09-13", "2026-2026", "13/09/26",
	} {
		if _, _, err := ParseDate(value); err == nil {
			t.Fatalf("%q: expected an error", value)
		}
	}
}

func TestParseTime(t *testing.T) {
	cases := []struct {
		value string
		start int
		end   int
	}{
		{"7 PM", 190000, 199999},
		{"7PM", 190000, 199999},
		{"07:30 PM", 193000, 193099},
		{"19:30", 193000, 193099},
		{"11:00", 110000, 110099},
		{"19:30:45", 193045, 193045},
		{"12 AM", 0, 9999},
		{"12 PM", 120000, 129999},
	}

	for _, c := range cases {
		start, end, err := ParseTime(c.value)
		if err != nil {
			t.Fatalf("%q: unexpected error %v", c.value, err)
		}
		if start != c.start || end != c.end {
			t.Fatalf("%q: expected [%d %d], got [%d %d]", c.value, c.start, c.end, start, end)
		}
	}
}

func TestParseTimeRejectsInvalid(t *testing.T) {
	for _, value := range []string{
		"", "noon", "25:00", "13 PM", "0 PM", "19:60", "19:30:60", "1:2:3:4",
	} {
		if _, _, err := ParseTime(value); err == nil {
			t.Fatalf("%q: expected an error", value)
		}
	}
}

func TestRecordValues(t *testing.T) {
	ts := time.Date(2026, 9, 13, 19, 30, 45, 0, time.UTC)

	if got := RecordDate(ts); got != 20260913 {
		t.Fatalf("expected 20260913, got %d", got)
	}

	if got := RecordTime(ts); got != 193045 {
		t.Fatalf("expected 193045, got %d", got)
	}
}
