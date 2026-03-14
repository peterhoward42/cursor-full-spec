package app

import "testing"

func TestFormatStoragePath_HappyPath(t *testing.T) {
	event := &TelemetryEvent{
		EventULID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		TimeUTC:   "2025-03-14T09:30:00Z",
	}

	got, err := FormatStoragePath(event)
	if err != nil {
		t.Fatalf("FormatStoragePath() error = %v, want nil", err)
	}

	want := "events/2025/03/14/09/01ARZ3NDEKTSV4RRFFQ69G5FAV"
	if got != want {
		t.Errorf("FormatStoragePath() = %q, want %q", got, want)
	}
}

func TestFormatStoragePath_NilEvent(t *testing.T) {
	t.Parallel()

	if _, err := FormatStoragePath(nil); err == nil {
		t.Errorf("FormatStoragePath(nil) error = nil, want non-nil")
	}
}

func TestFormatStoragePath_InvalidTime(t *testing.T) {
	t.Parallel()

	event := &TelemetryEvent{
		EventULID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		TimeUTC:   "not-a-time",
	}

	if _, err := FormatStoragePath(event); err == nil {
		t.Errorf("FormatStoragePath() error = nil for invalid time, want non-nil")
	}
}

