package app

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

// SerializeEventNDJSONGzip marshals the event as a single NDJSON line (JSON + newline)
// and returns the gzip-compressed bytes.
func SerializeEventNDJSONGzip(event *TelemetryEvent) ([]byte, error) {
	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}
	line, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("marshal event: %w", err)
	}
	line = append(line, '\n')

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(line); err != nil {
		_ = w.Close()
		return nil, fmt.Errorf("gzip write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("gzip close: %w", err)
	}
	return buf.Bytes(), nil
}

// ParseNDJSONGzip decompresses gzip and decodes the first NDJSON line into a TelemetryEvent.
// Used by tests to assert stored content; not used in production ingest path.
func ParseNDJSONGzip(data []byte) (*TelemetryEvent, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip reader: %w", err)
	}
	defer func() { _ = r.Close() }()
	dec := json.NewDecoder(r)
	var event TelemetryEvent
	if err := dec.Decode(&event); err != nil && err != io.EOF {
		return nil, fmt.Errorf("decode ndjson: %w", err)
	}
	return &event, nil
}
