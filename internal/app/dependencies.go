package app

// EventStorer stores telemetry events at a given path.
// StoreEventIfNotExists stores the event only if nothing exists at path; if something already exists at path, it does nothing and returns nil. It returns an error only on failure.
type EventStorer interface {
	StoreEventIfNotExists(path string, event *TelemetryEvent) error
}

// Dependencies holds interfaces to external systems.
// Start empty; add fields only when required by application behaviour.
type Dependencies struct {
	EventStorer EventStorer
}
