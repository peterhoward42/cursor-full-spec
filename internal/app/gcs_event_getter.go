package app

import (
	"context"
	"errors"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// GCSEventGetter implements EventGetter by reading all stored event objects from Google Cloud Storage.
// Objects are listed under the "events/" prefix and each object is expected to be NDJSON gzip as written by GCSEventStorer.
type GCSEventGetter struct {
	Bucket string
	Client *storage.Client
}

// GetAllStoredEvents lists all objects under the bucket's "events/" prefix, reads each object's content,
// and parses it as NDJSON gzip into a TelemetryEvent. Returns a slice of all events and nil error, or nil and an error.
func (g *GCSEventGetter) GetAllStoredEvents() ([]TelemetryEvent, error) {
	if g.Client == nil {
		return nil, errors.New("GCS client is nil")
	}
	if g.Bucket == "" {
		return nil, errors.New("GCS bucket name is empty")
	}

	ctx := context.Background()
	bucket := g.Client.Bucket(g.Bucket)
	it := bucket.Objects(ctx, &storage.Query{Prefix: "events/"})

	var events []TelemetryEvent
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		obj := bucket.Object(attrs.Name)
		r, err := obj.NewReader(ctx)
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(r)
		_ = r.Close()
		if err != nil {
			return nil, err
		}

		event, err := ParseNDJSONGzip(data)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}
