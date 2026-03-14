package app

import (
	"testing"
)

// eventWith returns a copy of validTelemetryEvent() with Event and optionally ProxyUserID overridden.
func eventWith(event string, proxyUserID string) TelemetryEvent {
	e := validTelemetryEvent()
	e.Event = event
	if proxyUserID != "" {
		e.ProxyUserID = proxyUserID
	}
	return e
}

func TestAnalyseEvents_EmptyInput(t *testing.T) {
	t.Parallel()

	got, err := AnalyseEvents(nil)
	if err != nil {
		t.Fatalf("AnalyseEvents(nil) err = %v, want nil", err)
	}
	assertReportZero(got, t)

	got, err = AnalyseEvents([]TelemetryEvent{})
	if err != nil {
		t.Fatalf("AnalyseEvents([]) err = %v, want nil", err)
	}
	assertReportZero(got, t)
}

func assertReportZero(r AnalysisReport, t *testing.T) {
	t.Helper()
	if r.HowManyPeopleHave.Launched != 0 || r.HowManyPeopleHave.LoadedAnExample != 0 ||
		r.HowManyPeopleHave.TriedToSignIn != 0 || r.HowManyPeopleHave.SucceededSigningIn != 0 ||
		r.HowManyPeopleHave.CreatedTheirOwnDrawing != 0 || r.HowManyPeopleHave.RetreivedTheirASavedDrawing != 0 ||
		r.TotalRecoverableErrors != 0 || r.TotalFatalErrors != 0 {
		t.Errorf("expected all-zero report, got %+v", r)
	}
}

func TestAnalyseEvents_CountsDistinctUsersPerPeopleEvent(t *testing.T) {
	t.Parallel()

	events := []TelemetryEvent{
		eventWith(EventLaunched, "user-a"),
		eventWith(EventLaunched, "user-a"),
		eventWith(EventLaunched, "user-b"),
	}

	got, err := AnalyseEvents(events)
	if err != nil {
		t.Fatalf("AnalyseEvents() err = %v, want nil", err)
	}
	if got.HowManyPeopleHave.Launched != 2 {
		t.Errorf("HowManyPeopleHave.Launched = %d, want 2 (distinct user-a, user-b)", got.HowManyPeopleHave.Launched)
	}
}

func TestAnalyseEvents_EachPeopleEventTypeCountedIndependently(t *testing.T) {
	t.Parallel()

	events := []TelemetryEvent{
		eventWith(EventLaunched, "user-1"),
		eventWith(EventLoadedExample, "user-1"),
		eventWith(EventSignInStarted, "user-1"),
		eventWith(EventSignInSuccess, "user-1"),
		eventWith(EventCreatedNewDrawing, "user-1"),
		eventWith(EventRetrievedSavedDrawing, "user-1"),
	}

	got, err := AnalyseEvents(events)
	if err != nil {
		t.Fatalf("AnalyseEvents() err = %v, want nil", err)
	}
	if got.HowManyPeopleHave.Launched != 1 {
		t.Errorf("Launched = %d, want 1", got.HowManyPeopleHave.Launched)
	}
	if got.HowManyPeopleHave.LoadedAnExample != 1 {
		t.Errorf("LoadedAnExample = %d, want 1", got.HowManyPeopleHave.LoadedAnExample)
	}
	if got.HowManyPeopleHave.TriedToSignIn != 1 {
		t.Errorf("TriedToSignIn = %d, want 1", got.HowManyPeopleHave.TriedToSignIn)
	}
	if got.HowManyPeopleHave.SucceededSigningIn != 1 {
		t.Errorf("SucceededSigningIn = %d, want 1", got.HowManyPeopleHave.SucceededSigningIn)
	}
	if got.HowManyPeopleHave.CreatedTheirOwnDrawing != 1 {
		t.Errorf("CreatedTheirOwnDrawing = %d, want 1", got.HowManyPeopleHave.CreatedTheirOwnDrawing)
	}
	if got.HowManyPeopleHave.RetreivedTheirASavedDrawing != 1 {
		t.Errorf("RetreivedTheirASavedDrawing = %d, want 1", got.HowManyPeopleHave.RetreivedTheirASavedDrawing)
	}
}

func TestAnalyseEvents_ErrorEventsCountTotalOccurrences(t *testing.T) {
	t.Parallel()

	events := []TelemetryEvent{
		eventWith(EventRecoverableJSError, "user-1"),
		eventWith(EventRecoverableJSError, "user-1"),
		eventWith(EventFatalJSError, "user-2"),
	}

	got, err := AnalyseEvents(events)
	if err != nil {
		t.Fatalf("AnalyseEvents() err = %v, want nil", err)
	}
	if got.TotalRecoverableErrors != 2 {
		t.Errorf("TotalRecoverableErrors = %d, want 2", got.TotalRecoverableErrors)
	}
	if got.TotalFatalErrors != 1 {
		t.Errorf("TotalFatalErrors = %d, want 1", got.TotalFatalErrors)
	}
}

func TestAnalyseEvents_UnknownEventTypesIgnored(t *testing.T) {
	t.Parallel()

	events := []TelemetryEvent{
		eventWith("unknown-event", "user-1"),
		eventWith(EventLaunched, "user-1"),
	}

	got, err := AnalyseEvents(events)
	if err != nil {
		t.Fatalf("AnalyseEvents() err = %v, want nil", err)
	}
	if got.HowManyPeopleHave.Launched != 1 {
		t.Errorf("Launched = %d, want 1 (unknown event should be ignored)", got.HowManyPeopleHave.Launched)
	}
}
