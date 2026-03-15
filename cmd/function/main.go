package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/pete/cursor-full-spec/internal/app"
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("storage client Close: %v", err)
		}
	}()

	bucket := os.Getenv("GCS_BUCKET")
	if bucket == "" {
		log.Fatal("GCS_BUCKET environment variable is not set")
	}

	deps := app.Dependencies{
		EventStorer: &app.GCSEventStorer{Bucket: bucket, Client: client},
		EventGetter: &app.GCSEventGetter{Bucket: bucket, Client: client},
	}
	a := app.NewApplication(deps)
	functions.HTTP("Function", a.ServeHTTP)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v", err)
	}
}
