package app

// FakeEventStorer is a test-double implementation of EventStorer that records calls.
// Use it in tests and as a placeholder at Dependencies construction sites.
// Set StoreErr to make StoreEventIfNotExists return that error (for failure-path tests).
type FakeEventStorer struct {
	// Stored holds each (path, data) passed to StoreEventIfNotExists for test assertions.
	Stored []StoredCall
	// StoreErr, when non-nil, is returned from StoreEventIfNotExists without recording the call.
	StoreErr error
}

// StoredCall represents a single call to StoreEventIfNotExists.
type StoredCall struct {
	Path string
	Data []byte
}

// StoreEventIfNotExists records the path and data only if path is not already in Stored; otherwise does nothing.
// Returns f.StoreErr if set; otherwise returns nil.
func (f *FakeEventStorer) StoreEventIfNotExists(path string, data []byte) error {
	if f.StoreErr != nil {
		return f.StoreErr
	}
	for _, c := range f.Stored {
		if c.Path == path {
			return nil
		}
	}
	// Copy so caller cannot mutate stored data.
	dup := make([]byte, len(data))
	copy(dup, data)
	f.Stored = append(f.Stored, StoredCall{Path: path, Data: dup})
	return nil
}
