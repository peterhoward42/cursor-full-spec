package app

// AnalysisReport holds aggregate counts derived from TelemetryEvents.
// HowManyPeopleHave fields count distinct ProxyUserID per event type;
// error totals count all events of that type.
type AnalysisReport struct {
	HowManyPeopleHave struct {
		Launched                   int `json:"Launched"`
		LoadedAnExample            int `json:"LoadedAnExample"`
		TriedToSignIn              int `json:"TriedToSignIn"`
		SucceededSigningIn        int `json:"SucceededSigningIn"`
		CreatedTheirOwnDrawing    int `json:"CreatedTheirOwnDrawing"`
		RetreivedTheirASavedDrawing int `json:"RetreivedTheirASavedDrawing"`
	} `json:"HowManyPeopleHave"`
	TotalRecoverableErrors int `json:"TotalRecoverableErrors"`
	TotalFatalErrors       int `json:"TotalFatalErrors"`
}

// AnalyseEvents consumes a slice of TelemetryEvents and produces an AnalysisReport.
// For "people" event types (e.g. launched, loaded-example), counts distinct ProxyUserID.
// For recoverable/fatal JS error events, counts total occurrences.
func AnalyseEvents(events []TelemetryEvent) (AnalysisReport, error) {
	var report AnalysisReport
	if len(events) == 0 {
		return report, nil
	}

	// Distinct ProxyUserID per "people" event type
	launched := make(map[string]struct{})
	loadedExample := make(map[string]struct{})
	signInStarted := make(map[string]struct{})
	signInSuccess := make(map[string]struct{})
	createdNewDrawing := make(map[string]struct{})
	retrievedSavedDrawing := make(map[string]struct{})

	for _, e := range events {
		uid := e.ProxyUserID
		switch e.Event {
		case EventLaunched:
			launched[uid] = struct{}{}
		case EventLoadedExample:
			loadedExample[uid] = struct{}{}
		case EventSignInStarted:
			signInStarted[uid] = struct{}{}
		case EventSignInSuccess:
			signInSuccess[uid] = struct{}{}
		case EventCreatedNewDrawing:
			createdNewDrawing[uid] = struct{}{}
		case EventRetrievedSavedDrawing:
			retrievedSavedDrawing[uid] = struct{}{}
		case EventRecoverableJSError:
			report.TotalRecoverableErrors++
		case EventFatalJSError:
			report.TotalFatalErrors++
		}
	}

	report.HowManyPeopleHave.Launched = len(launched)
	report.HowManyPeopleHave.LoadedAnExample = len(loadedExample)
	report.HowManyPeopleHave.TriedToSignIn = len(signInStarted)
	report.HowManyPeopleHave.SucceededSigningIn = len(signInSuccess)
	report.HowManyPeopleHave.CreatedTheirOwnDrawing = len(createdNewDrawing)
	report.HowManyPeopleHave.RetreivedTheirASavedDrawing = len(retrievedSavedDrawing)

	return report, nil
}
