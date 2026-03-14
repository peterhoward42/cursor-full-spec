package app

import (
	"io"
	"net/http"
)

// Application owns request-handling logic and collaborates via Dependencies.
type Application struct {
	deps Dependencies
}

// NewApplication constructs an Application with explicit Dependencies.
func NewApplication(deps Dependencies) *Application {
	return &Application{deps: deps}
}

// ServeHTTP dispatches by method to the appropriate Application method.
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.IngestEvent(w, r)
	case http.MethodGet:
		a.AnalysisReport(w, r)
	case http.MethodOptions:
		a.CORS(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// IngestEvent handles POST of a single TelemetryEvent: parse and validate payload,
// serialise as NDJSON gzip, then store at the path from FormatStoragePath.
func (a *Application) IngestEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}
	event, err := ParseAndValidateTelemetryEvent(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path, err := FormatStoragePath(event)
	if err != nil {
		http.Error(w, "storage path", http.StatusInternalServerError)
		return
	}
	data, err := SerializeEventNDJSONGzip(event)
	if err != nil {
		http.Error(w, "serialise event", http.StatusInternalServerError)
		return
	}
	if err := a.deps.EventStorer.StoreEventIfNotExists(path, data); err != nil {
		http.Error(w, "store event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// AnalysisReport handles GET and returns an AnalysisReport (placeholder).
func (a *Application) AnalysisReport(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// CORS handles OPTIONS with a CORS response for browser requests (placeholder).
func (a *Application) CORS(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
