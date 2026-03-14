package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplication_ServeHTTP_GET_returnsNotImplemented(t *testing.T) {
	// Given: Application with empty Dependencies (placeholder behaviour).
	deps := Dependencies{
		EventStorer: &FakeEventStorer{},
	}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// When: GET request is handled.
	app.ServeHTTP(rec, req)

	// Then: placeholder responds with 501 (derived from "no implementation yet").
	if rec.Code != http.StatusNotImplemented {
		t.Errorf("ServeHTTP(GET) status = %d, want %d", rec.Code, http.StatusNotImplemented)
	}
}
