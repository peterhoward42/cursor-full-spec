package app

// EventStorer stores serialised telemetry event data at a given path.
// StoreEventIfNotExists writes data only if nothing exists at path; if something already exists at path, it does nothing and returns nil. It returns an error only on failure.
type EventStorer interface {
	StoreEventIfNotExists(path string, data []byte) error
}

// EventGetter provides access to all stored TelemetryEvents.
type EventGetter interface {
	GetAllStoredEvents() ([]TelemetryEvent, error)
}

// Dependencies holds interfaces to external systems.
// Start empty; add fields only when required by application behaviour.
type Dependencies struct {
	EventStorer  EventStorer
	EventGetter  EventGetter
}
