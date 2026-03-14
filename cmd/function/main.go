package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/pete/cursor-full-spec/internal/app"
)

func main() {
	deps := app.Dependencies{
		EventStorer: &app.FakeEventStorer{},
		EventGetter: &app.FakeEventGetter{},
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
