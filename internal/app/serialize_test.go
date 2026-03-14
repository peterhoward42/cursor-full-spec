package app

import (
	"testing"
)

func TestSerializeEventNDJSONGzip_NilEvent_ReturnsError(t *testing.T) {
	_, err := SerializeEventNDJSONGzip(nil)
	if err == nil {
		t.Fatal("SerializeEventNDJSONGzip(nil) error = nil, want non-nil")
	}
}

func TestSerializeEventNDJSONGzip_RoundTrip(t *testing.T) {
	event := TelemetryEvent{
		SchemaVersion: 1,
		EventULID:     "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		ProxyUserID:   "550e8400-e29b-41d4-a716-446655440000",
		TimeUTC:       "2025-03-14T12:00:00Z",
		Visit:         1,
		Event:         EventLaunched,
		Parameters:    "",
	}
	data, err := SerializeEventNDJSONGzip(&event)
	if err != nil {
		t.Fatalf("SerializeEventNDJSONGzip() error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("SerializeEventNDJSONGzip() returned empty data")
	}
	decoded, err := ParseNDJSONGzip(data)
	if err != nil {
		t.Fatalf("ParseNDJSONGzip() error = %v", err)
	}
	if *decoded != event {
		t.Errorf("round-trip event = %+v, want %+v", *decoded, event)
	}
}
