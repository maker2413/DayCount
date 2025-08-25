package compute

import (
	"testing"
	"time"
)

func mustDate(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestDaysBetween_Basic(t *testing.T) {
	a := mustDate(2024, 1, 1)
	b := mustDate(2024, 1, 11)
	if got := DaysBetween(a, b); got != 10 {
		t.Fatalf("expected 10 got %d", got)
	}
}

func TestDaysBetween_SameDay(t *testing.T) {
	a := mustDate(2024, 5, 5)
	if got := DaysBetween(a, a); got != 0 {
		t.Fatalf("expected 0 got %d", got)
	}
}

func TestDaysSince_FutureNegative(t *testing.T) {
	now := mustDate(2024, 6, 1)
	future := mustDate(2024, 6, 10)
	got := DaysSince(future, now)
	if got >= 0 {
		t.Fatalf("expected negative for future date got %d", got)
	}
	if got != -9 { // days between 1->10 is 9
		t.Fatalf("expected -9 got %d", got)
	}
}

func TestDaysBetween_LeapSpan(t *testing.T) {
	start := mustDate(2019, 2, 28)
	end := mustDate(2019, 3, 1)
	if got := DaysBetween(start, end); got != 1 {
		t.Fatalf("expected 1 got %d", got)
	}

	// Across leap day year
	start = mustDate(2020, 2, 28)
	end = mustDate(2020, 3, 1)
	if got := DaysBetween(start, end); got != 2 {
		t.Fatalf("expected 2 across leap day got %d", got)
	}
}

func TestDaysSinceNow_NonZero(t *testing.T) {
	// Just ensure it runs; cannot assert exact value w/out controlling now.
	d := time.Now().UTC().Add(-48 * time.Hour)
	if got := DaysSinceNow(d); got < 1 {
		t.Fatalf("expected at least 1 got %d", got)
	}
}
