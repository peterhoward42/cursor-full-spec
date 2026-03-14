package app

import (
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

// IngestEvent handles POST of a single TelemetryEvent (placeholder).
func (a *Application) IngestEvent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// AnalysisReport handles GET and returns an AnalysisReport (placeholder).
func (a *Application) AnalysisReport(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// CORS handles OPTIONS with a CORS response for browser requests (placeholder).
func (a *Application) CORS(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
