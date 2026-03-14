package app

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseAndValidateTelemetryEvent_HappyPath(t *testing.T) {
	t.Parallel()

	event := TelemetryEvent{
		SchemaVersion: 1,
		EventULID:     "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		ProxyUserID:   "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		TimeUTC:       "2025-03-14T09:30:00Z",
		Visit:         42,
		Event:         EventLaunched,
		Parameters:    "example-parameters",
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v, want nil", err)
	}

	got, err := ParseAndValidateTelemetryEvent(string(payloadBytes))
	if err != nil {
		t.Fatalf("ParseAndValidateTelemetryEvent() error = %v, want nil", err)
	}

	if *got != event {
		t.Errorf("ParseAndValidateTelemetryEvent() = %+v, want %+v", *got, event)
	}
}

func TestParseAndValidateTelemetryEvent_EmptyPayload(t *testing.T) {
	t.Parallel()

	if _, err := ParseAndValidateTelemetryEvent("   "); err == nil {
		t.Fatalf("ParseAndValidateTelemetryEvent() error = nil, want non-nil")
	}
}

func TestParseAndValidateTelemetryEvent_InvalidJSON(t *testing.T) {
	t.Parallel()

	if _, err := ParseAndValidateTelemetryEvent("{invalid-json"); err == nil {
		t.Fatalf("ParseAndValidateTelemetryEvent() error = nil for invalid JSON, want non-nil")
	}
}

func TestParseAndValidateTelemetryEvent_UnknownField(t *testing.T) {
	t.Parallel()

	payload := `{"SchemaVersion":1,"EventULID":"01ARZ3NDEKTSV4RRFFQ69G5FAV","ProxyUserID":"123e4567-e89b-12d3-a456-426614174000","TimeUTC":"2025-03-14T09:30:00Z","Visit":1,"Event":"launched","Parameters":"ok","Unexpected":"value"}`

	if _, err := ParseAndValidateTelemetryEvent(payload); err == nil {
		t.Fatalf("ParseAndValidateTelemetryEvent() error = nil for payload with unknown field, want non-nil")
	}
}

func TestParseAndValidateTelemetryEvent_FieldValidations(t *testing.T) {
	t.Parallel()

	base := TelemetryEvent{
		SchemaVersion: 1,
		EventULID:     "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		ProxyUserID:   "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		TimeUTC:       "2025-03-14T09:30:00Z",
		Visit:         1,
		Event:         EventLaunched,
		Parameters:    "ok",
	}

	cases := []struct {
		name    string
		mutate  func(e *TelemetryEvent)
		wantErr bool
	}{
		{
			name: "schema version too high",
			mutate: func(e *TelemetryEvent) {
				e.SchemaVersion = 2
			},
			wantErr: true,
		},
		{
			name: "empty ULID",
			mutate: func(e *TelemetryEvent) {
				e.EventULID = ""
			},
			wantErr: true,
		},
		{
			name: "short ULID",
			mutate: func(e *TelemetryEvent) {
				e.EventULID = "too-short"
			},
			wantErr: true,
		},
		{
			name: "empty proxy user id",
			mutate: func(e *TelemetryEvent) {
				e.ProxyUserID = ""
			},
			wantErr: true,
		},
		{
			name: "invalid proxy user id format",
			mutate: func(e *TelemetryEvent) {
				e.ProxyUserID = "not-a-uuid"
			},
			wantErr: true,
		},
		{
			name: "invalid proxy user id version",
			mutate: func(e *TelemetryEvent) {
				// A version 1 UUID.
				e.ProxyUserID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
			},
			wantErr: true,
		},
		{
			name: "empty time",
			mutate: func(e *TelemetryEvent) {
				e.TimeUTC = ""
			},
			wantErr: true,
		},
		{
			name: "bad time format",
			mutate: func(e *TelemetryEvent) {
				e.TimeUTC = "not-a-time"
			},
			wantErr: true,
		},
		{
			name: "visit too low",
			mutate: func(e *TelemetryEvent) {
				e.Visit = 0
			},
			wantErr: true,
		},
		{
			name: "visit too high",
			mutate: func(e *TelemetryEvent) {
				e.Visit = 100001
			},
			wantErr: true,
		},
		{
			name: "event name too short",
			mutate: func(e *TelemetryEvent) {
				e.Event = "abc"
			},
			wantErr: true,
		},
		{
			name: "event name too long",
			mutate: func(e *TelemetryEvent) {
				e.Event = fmt.Sprintf("%041s", "a")
			},
			wantErr: true,
		},
		{
			name: "event name unknown",
			mutate: func(e *TelemetryEvent) {
				e.Event = "unknown-event"
			},
			wantErr: true,
		},
		{
			name: "parameters too long",
			mutate: func(e *TelemetryEvent) {
				e.Parameters = fmt.Sprintf("%081s", "a")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ev := base
			tc.mutate(&ev)

			payloadBytes, err := json.Marshal(ev)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v, want nil", err)
			}

			_, err = ParseAndValidateTelemetryEvent(string(payloadBytes))
			if (err != nil) != tc.wantErr {
				t.Fatalf("ParseAndValidateTelemetryEvent() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

