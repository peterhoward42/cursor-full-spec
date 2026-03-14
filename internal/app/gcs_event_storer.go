package app

import (
	"context"
	"encoding/json"
	"errors"

	"cloud.google.com/go/storage"
)

// GCSEventStorer implements EventStorer by writing events as JSON objects to Google Cloud Storage.
// StoreEventIfNotExists writes only when no object exists at path; otherwise it does nothing and returns nil.
type GCSEventStorer struct {
	Bucket string
	Client *storage.Client
}

// StoreEventIfNotExists writes the event to path in GCS only if no object exists there.
// If an object is already present at path, it does nothing and returns nil.
func (g *GCSEventStorer) StoreEventIfNotExists(path string, event *TelemetryEvent) error {
	if g.Client == nil {
		return errors.New("GCS client is nil")
	}
	if g.Bucket == "" {
		return errors.New("GCS bucket name is empty")
	}
	if event == nil {
		return errors.New("event is nil")
	}

	ctx := context.Background()
	obj := g.Client.Bucket(g.Bucket).Object(path)

	_, err := obj.Attrs(ctx)
	if err == nil {
		return nil
	}
	if !errors.Is(err, storage.ErrObjectNotExist) {
		return err
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	w := obj.NewWriter(ctx)
	if _, err := w.Write(body); err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}
