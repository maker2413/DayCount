package parse

import (
	"testing"
	"time"
)

func TestParseDate_SupportedFormats(t *testing.T) {
	cases := []string{
		"2024-12-31",
		"2024/12/31",
		"12/31/2024",
		"2024-12-31T23:59:59Z",
		"2024-12-31T23:59:59",
		"2024-12-31 23:59:59",
	}
	for _, in := range cases {
		got, err := ParseDate(in)
		if err != nil {
			// use t.Fatalf to stop immediate iteration for clarity
			t.Fatalf("ParseDate(%q) unexpected error: %v", in, err)
		}
		if got.Location() != time.UTC {
			t.Errorf("expected UTC location for %q", in)
		}
		if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 {
			t.Errorf("expected midnight normalization for %q got %v", in, got)
		}
	}
}

func TestParseDate_LeapDay(t *testing.T) {
	got, err := ParseDate("2020-02-29T12:30:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if got.Year() != 2020 || got.Month() != 2 || got.Day() != 29 {
		t.Errorf("bad leap day parse: %v", got)
	}
}

func TestParseDate_Error(t *testing.T) {
	_, err := ParseDate("not-a-date")
	if err == nil {
		t.Fatal("expected error")
	}
	if err != ErrUnrecognizedFormat {
		t.Fatalf("expected ErrUnrecognizedFormat got %v", err)
	}
}

func TestParseDate_StripsTime(t *testing.T) {
	got, err := ParseDate("2024-07-04T05:06:07Z")
	if err != nil {
		t.Fatal(err)
	}
	if got.Hour() != 0 {
		t.Errorf("expected hour 0 got %d", got.Hour())
	}
}
