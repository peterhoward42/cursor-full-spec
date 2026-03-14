package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ParseAndValidateTelemetryEvent parses the given JSON payload into a TelemetryEvent
// and performs comprehensive validation of its fields.
func ParseAndValidateTelemetryEvent(payload string) (*TelemetryEvent, error) {
	if strings.TrimSpace(payload) == "" {
		return nil, fmt.Errorf("payload is empty")
	}

	var event TelemetryEvent

	dec := json.NewDecoder(strings.NewReader(payload))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&event); err != nil {
		return nil, fmt.Errorf("decode telemetry event: %w", err)
	}

	if err := validateTelemetryEvent(&event); err != nil {
		return nil, err
	}

	return &event, nil
}

func validateTelemetryEvent(event *TelemetryEvent) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	if event.SchemaVersion < 0 || event.SchemaVersion > 1 {
		return fmt.Errorf("invalid SchemaVersion: must be between 0 and 1 inclusive")
	}

	if event.EventULID == "" {
		return fmt.Errorf("invalid EventULID: must not be empty")
	}
	if len(event.EventULID) != 26 {
		return fmt.Errorf("invalid EventULID: must be 26 characters long")
	}

	if event.ProxyUserID == "" {
		return fmt.Errorf("invalid ProxyUserID: must not be empty")
	}
	u, err := uuid.Parse(event.ProxyUserID)
	if err != nil {
		return fmt.Errorf("invalid ProxyUserID: %w", err)
	}
	if u.Version() != 4 {
		return fmt.Errorf("invalid ProxyUserID: must be a UUIDv4")
	}

	if event.TimeUTC == "" {
		return fmt.Errorf("invalid TimeUTC: must not be empty")
	}
	if _, err := time.Parse(time.RFC3339, event.TimeUTC); err != nil {
		return fmt.Errorf("invalid TimeUTC: %w", err)
	}

	if event.Visit < 1 || event.Visit > 100000 {
		return fmt.Errorf("invalid Visit: must be between 1 and 100000 inclusive")
	}

	if l := len(event.Event); l < 4 || l > 40 {
		return fmt.Errorf("invalid Event: length must be between 4 and 40 characters")
	}
	if !isKnownEventName(event.Event) {
		return fmt.Errorf("invalid Event: %q is not a recognised event name", event.Event)
	}

	if len(event.Parameters) > 80 {
		return fmt.Errorf("invalid Parameters: length must be at most 80 characters")
	}

	return nil
}

func isKnownEventName(name string) bool {
	switch name {
	case EventLaunched,
		EventLoadedExample,
		EventSignInStarted,
		EventSignInSuccess,
		EventCreatedNewDrawing,
		EventRetrievedSavedDrawing,
		EventRecoverableJSError,
		EventFatalJSError:
		return true
	default:
		return false
	}
}

