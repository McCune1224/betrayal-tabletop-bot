package util

import (
	"time"
)

// Helper to get current time in EST
func GetEstTimeStamp() string {
	// get current time in est
	est := time.Now().UTC().Add(-5 * time.Hour)

	// format similar to Oct 24 11:00:00
	// without the est at the end
	return est.Format("Jan 2 15:04:05")
}

func GetEstTime(t time.Time) string {
  est := t.UTC().Add(-5 * time.Hour)
  return est.Format("Jan 2 15:04:05")
}

func GetEstTimeStampFromDuration(d time.Duration) string {
	// get current time in est
	est := time.Now().UTC().Add(-5 * time.Hour)
	est = est.Add(d)
	// format similar to Oct 24 11:00:00
	// without the est at the end
	return est.Format("Jan 2 15:04:05")
}
