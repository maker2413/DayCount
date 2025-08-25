package parse

import (
	"errors"
	"strings"
	"time"
)

var (
	// ErrUnrecognizedFormat is returned when the input cannot be parsed as a supported date.
	ErrUnrecognizedFormat = errors.New("unrecognized date format")
)

// supported layouts tried in order. Comment each for clarity.
var layouts = []string{
	time.DateOnly,         // 2006-01-02
	"2006/01/02",          // 2006/01/02
	"01/02/2006",          // 01/02/2006 (US)
	time.RFC3339,          // 2006-01-02T15:04:05Z07:00
	"2006-01-02T15:04:05", // no zone
	"2006-01-02 15:04:05", // space separator, no zone
}

// ParseDate parses many common date (or datetime) formats into a UTC midnight time.Time.
// Time portion (if any) is discarded. If no time zone provided, local time is assumed before
// normalizing to date components in UTC.
func ParseDate(input string) (time.Time, error) {
	s := strings.TrimSpace(input)
	var parsed time.Time
	var err error
	for _, layout := range layouts {
		parsed, err = time.Parse(layout, s)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, ErrUnrecognizedFormat
	}

	// Normalize to date components in its own location (already in location per layout parsing)
	y, m, d := parsed.Date()

	// Reconstruct at midnight UTC to avoid DST irregularities.
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
}
