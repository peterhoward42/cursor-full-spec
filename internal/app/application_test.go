package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplication_ServeHTTP_GET_returnsReport(t *testing.T) {
	t.Parallel()
	deps := Dependencies{
		EventStorer: &FakeEventStorer{},
		EventGetter: &FakeEventGetter{},
	}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("ServeHTTP(GET) status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
	}
	var report AnalysisReport
	if err := json.Unmarshal(rec.Body.Bytes(), &report); err != nil {
		t.Fatalf("response body not valid JSON: %v", err)
	}
	assertReportZero(report, t)
}

func TestApplication_AnalysisReport_EventGetterError_Returns500(t *testing.T) {
	t.Parallel()
	deps := Dependencies{
		EventStorer: &FakeEventStorer{},
		EventGetter: &FakeEventGetter{GetErr: errors.New("get failed")},
	}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	app.AnalysisReport(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("AnalysisReport(EventGetter error) status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestApplication_AnalysisReport_EmptyEvents_Returns200AndZeroReport(t *testing.T) {
	t.Parallel()
	getter := &FakeEventGetter{Events: []TelemetryEvent{}}
	deps := Dependencies{EventStorer: &FakeEventStorer{}, EventGetter: getter}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	app.AnalysisReport(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("AnalysisReport(empty events) status = %d, want %d", rec.Code, http.StatusOK)
	}
	var report AnalysisReport
	if err := json.Unmarshal(rec.Body.Bytes(), &report); err != nil {
		t.Fatalf("response body not valid JSON: %v", err)
	}
	assertReportZero(report, t)
}

func TestApplication_AnalysisReport_WithEvents_Returns200AndReportFromAnalyser(t *testing.T) {
	t.Parallel()
	events := []TelemetryEvent{
		eventWith(EventLaunched, "user-a"),
		eventWith(EventLaunched, "user-b"),
		eventWith(EventRecoverableJSError, ""),
	}
	getter := &FakeEventGetter{Events: events}
	deps := Dependencies{EventStorer: &FakeEventStorer{}, EventGetter: getter}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	app.AnalysisReport(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("AnalysisReport(events) status = %d, want %d", rec.Code, http.StatusOK)
	}
	var report AnalysisReport
	if err := json.Unmarshal(rec.Body.Bytes(), &report); err != nil {
		t.Fatalf("response body not valid JSON: %v", err)
	}
	// Orchestration: EventGetter events were passed to AnalyseEvents and response matches.
	if report.HowManyPeopleHave.Launched != 2 {
		t.Errorf("HowManyPeopleHave.Launched = %d, want 2", report.HowManyPeopleHave.Launched)
	}
	if report.TotalRecoverableErrors != 1 {
		t.Errorf("TotalRecoverableErrors = %d, want 1", report.TotalRecoverableErrors)
	}
}

func TestApplication_IngestEvent_HappyPath(t *testing.T) {
	t.Parallel()
	// Given: valid payload and a fake storer.
	event := validTelemetryEvent()
	payload := mustMarshalEvent(t, event)
	expectedPath, err := FormatStoragePath(&event)
	if err != nil {
		t.Fatalf("FormatStoragePath() error = %v", err)
	}
	storer := &FakeEventStorer{}
	deps := Dependencies{EventStorer: storer, EventGetter: &FakeEventGetter{}}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	// When: POST with valid body is handled.
	app.IngestEvent(rec, req)

	// Then: 201 and storer received path from FormatStoragePath and serialised event data.
	if rec.Code != http.StatusCreated {
		t.Errorf("IngestEvent() status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if len(storer.Stored) != 1 {
		t.Fatalf("storer.Stored length = %d, want 1", len(storer.Stored))
	}
	if storer.Stored[0].Path != expectedPath {
		t.Errorf("storer.Stored[0].Path = %q, want %q", storer.Stored[0].Path, expectedPath)
	}
	decoded, err := ParseNDJSONGzip(storer.Stored[0].Data)
	if err != nil {
		t.Fatalf("ParseNDJSONGzip(stored data) error = %v", err)
	}
	if *decoded != event {
		t.Errorf("stored event = %+v, want %+v", *decoded, event)
	}
}

func TestApplication_IngestEvent_InvalidPayload_Returns400(t *testing.T) {
	t.Parallel()
	storer := &FakeEventStorer{}
	deps := Dependencies{EventStorer: storer, EventGetter: &FakeEventGetter{}}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("   "))
	rec := httptest.NewRecorder()

	app.IngestEvent(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("IngestEvent(empty body) status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if len(storer.Stored) != 0 {
		t.Errorf("storer.Stored length = %d, want 0 (storer not called)", len(storer.Stored))
	}
}

func TestApplication_IngestEvent_ValidationError_Returns400(t *testing.T) {
	t.Parallel()
	storer := &FakeEventStorer{}
	deps := Dependencies{EventStorer: storer, EventGetter: &FakeEventGetter{}}
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"invalid": true}`))
	rec := httptest.NewRecorder()

	app.IngestEvent(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("IngestEvent(invalid JSON) status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if len(storer.Stored) != 0 {
		t.Errorf("storer.Stored length = %d, want 0", len(storer.Stored))
	}
}

func TestApplication_IngestEvent_StorerError_Returns500(t *testing.T) {
	t.Parallel()
	storer := &FakeEventStorer{StoreErr: errors.New("storage unavailable")}
	deps := Dependencies{EventStorer: storer, EventGetter: &FakeEventGetter{}}
	event := validTelemetryEvent()
	payload := mustMarshalEvent(t, event)
	app := NewApplication(deps)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	app.IngestEvent(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("IngestEvent(storer error) status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if len(storer.Stored) != 0 {
		t.Errorf("storer.Stored length = %d, want 0 (call not recorded when error)", len(storer.Stored))
	}
}

func TestApplication_IngestEvent_IdempotentWhenAlreadyStored(t *testing.T) {
	t.Parallel()
	event := validTelemetryEvent()
	payload := mustMarshalEvent(t, event)
	storer := &FakeEventStorer{}
	deps := Dependencies{EventStorer: storer, EventGetter: &FakeEventGetter{}}
	app := NewApplication(deps)

	// First POST: stored.
	rec1 := httptest.NewRecorder()
	app.IngestEvent(rec1, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload)))
	if rec1.Code != http.StatusCreated {
		t.Errorf("first POST status = %d, want %d", rec1.Code, http.StatusCreated)
	}
	path := storer.Stored[0].Path

	// Pre-seed same path so second store is "already exists".
	storer.Stored = append(storer.Stored, StoredCall{Path: path, Data: storer.Stored[0].Data})

	// Second POST: StoreEventIfNotExists does nothing, returns nil; still 201.
	rec2 := httptest.NewRecorder()
	app.IngestEvent(rec2, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload)))
	if rec2.Code != http.StatusCreated {
		t.Errorf("second POST status = %d, want %d", rec2.Code, http.StatusCreated)
	}
	// Still only one distinct path stored (Fake stores one then skips duplicate path).
	if len(storer.Stored) != 2 {
		t.Errorf("storer.Stored length = %d, want 2 (one real, one pre-seed)", len(storer.Stored))
	}
}
