package app

import (
	"encoding/json"
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

// AnalysisReport handles GET: fetches all stored events via EventGetter, builds an AnalysisReport with AnalyseEvents, and returns it as JSON.
func (a *Application) AnalysisReport(w http.ResponseWriter, _ *http.Request) {
	events, err := a.deps.EventGetter.GetAllStoredEvents()
	if err != nil {
		http.Error(w, "get events", http.StatusInternalServerError)
		return
	}
	report, err := AnalyseEvents(events)
	if err != nil {
		http.Error(w, "analysis", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(report)
	if err != nil {
		http.Error(w, "encode", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// CORS handles OPTIONS preflight: responds with CORS headers and 204 No Content.
// CORS headers are sent only for OPTIONS, not for GET or POST.
func (a *Application) CORS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusNoContent)
}
