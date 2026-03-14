# Create an encapsulated report analysiser

## Context

- The document /docs/specs/overview-context outlines an analysis process that consumes a
set of TelemetryEvent(s) in order to produce an AnalysisReport
- This code generation step will provide an encapsulated Analyser to do that job.

## Instructions

- Write a function with this signature

```
AnalyseEvents(events []TelemetryEvent) (AnalysisReport, error)
```

- Start by attempting to deduce the analysis logic yourself.
- If you have insufficient information, stop and consult me
- If you are able to proceed then produce tests for the analyser
- Do not wire the analyser up yet.