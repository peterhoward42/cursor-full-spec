package app

// FakeEventStorer is a test-double implementation of EventStorer that records calls
// and never returns an error. Use it in tests and as a placeholder at Dependencies construction sites.
type FakeEventStorer struct {
	// Stored holds each (path, event) passed to StoreEventIfNotExists for test assertions.
	Stored []StoredCall
}

// StoredCall represents a single call to StoreEventIfNotExists.
type StoredCall struct {
	Path  string
	Event *TelemetryEvent
}

// StoreEventIfNotExists records the path and event only if path is not already in Stored; otherwise does nothing. Returns nil.
func (f *FakeEventStorer) StoreEventIfNotExists(path string, event *TelemetryEvent) error {
	for _, c := range f.Stored {
		if c.Path == path {
			return nil
		}
	}
	f.Stored = append(f.Stored, StoredCall{Path: path, Event: event})
	return nil
}
