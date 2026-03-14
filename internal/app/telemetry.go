package app

import (
	"errors"
	"fmt"
	"time"
)

// TelemetryEvent represents a single telemetry event emitted by the web app.
// Its fields and validation constraints are defined by the overview-context specification.
type TelemetryEvent struct {
	SchemaVersion int
	EventULID     string
	ProxyUserID   string
	TimeUTC       string
	Visit         int
	Event         string
	Parameters    string
}

// FormatStoragePath returns the slash-delimited storage path for the given TelemetryEvent.
//
// The format is:
//   events/year/month/day/hour/<ulid>
//
// where the year, month, day and hour components are derived from the event's TimeUTC
// field (RFC3339, in UTC) and the ULID comes from EventULID.
func FormatStoragePath(event *TelemetryEvent) (string, error) {
	if event == nil {
		return "", errors.New("event is nil")
	}

	if event.EventULID == "" {
		return "", errors.New("event ULID is empty")
	}

	if event.TimeUTC == "" {
		return "", errors.New("event time is empty")
	}

	t, err := time.Parse(time.RFC3339, event.TimeUTC)
	if err != nil {
		return "", fmt.Errorf("parse event time: %w", err)
	}

	t = t.UTC()
	year, month, day := t.Date()
	hour := t.Hour()

	path := fmt.Sprintf("events/%04d/%02d/%02d/%02d/%s", year, int(month), day, hour, event.EventULID)
	return path, nil
}

