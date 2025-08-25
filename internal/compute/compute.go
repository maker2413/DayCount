package compute

import "time"

// DaysBetween returns integer whole days (24h blocks) between from and to
// as to - from. Result negative if to is before from. Times are normalized
// to their date components in UTC before calculation.
func DaysBetween(from, to time.Time) int {
	fy, fm, fd := from.Date()
	ty, tm, td := to.Date()
	f := time.Date(fy, fm, fd, 0, 0, 0, 0, time.UTC)
	t := time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	dur := t.Sub(f)
	return int(dur.Hours() / 24)
}

// DaysSince returns days from the given past (or future) date to 'now'.
// Negative if date is in the future relative to now.
func DaysSince(date time.Time, now time.Time) int {
	return DaysBetween(date, now)
}

// DaysSinceNow convenience using time.Now().UTC().
func DaysSinceNow(date time.Time) int {
	return DaysSince(date, time.Now().UTC())
}
