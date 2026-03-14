package app

// EventStorer stores telemetry events at a given path.
// StoreEventIfNotExists stores the event only if nothing exists at path; it returns an error on conflict or failure.
type EventStorer interface {
	StoreEventIfNotExists(path string, event *TelemetryEvent) error
}

// Dependencies holds interfaces to external systems.
// Start empty; add fields only when required by application behaviour.
type Dependencies struct {
	EventStorer EventStorer
}
