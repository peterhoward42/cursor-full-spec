package app

// FakeEventGetter is a test-double implementation of EventGetter.
// Use it in tests and as a placeholder at Dependencies construction sites.
// Set Events to control the slice returned from GetAllStoredEvents; set GetErr to simulate errors.
type FakeEventGetter struct {
	// Events is the slice returned by GetAllStoredEvents when GetErr is nil.
	Events []TelemetryEvent
	// GetErr, when non-nil, is returned from GetAllStoredEvents without returning Events.
	GetErr error
}

// GetAllStoredEvents returns a copy of f.Events, or (nil, f.GetErr) if GetErr is set.
func (f *FakeEventGetter) GetAllStoredEvents() ([]TelemetryEvent, error) {
	if f.GetErr != nil {
		return nil, f.GetErr
	}
	out := make([]TelemetryEvent, len(f.Events))
	copy(out, f.Events)
	return out, nil
}
