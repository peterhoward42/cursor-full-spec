# This is an overview specification for a complete software system

- The specification will be used to generate the code required using AI
- The specification defines an overview of the system required, to provide high level context for you
- You must not start generating code from this specification alone because we will do that incrementally in a set of smaller steps

## Key Data Types

- The required system behaviour will be defined in terms of operations of a TelemetryEvent type and a AnalysisReport type.

## High level architecture

- The system is a Google Cloud Function implemented with the Go coding language and the
  Google Functions Framework package, using the "source deployment" option.

## High level behaviour description

### Event Ingestion 

- Handle a POST request that receives a JSON encoded single TelemetryEvent payload
- Each event should be stored in GoogleCloudStorage in a bucket named "events" as NDJSON, gzip compressed.
- The slash-delimited storage path should encode the time bucket using the time field inside
  the Event, and the event's ULID as shown by the following example.

```
events/year/month/day/hour/<ulid>
```

- The POST request should do nothing if the event already exists

### Analysis report generation

- Handle a GET request that returns a JSON encoded AnalysisReport
- Derive the analysis report by reading and reasoning over all the stored events

## The Telemetry Event data type

```
type TelemetryEvent struct {
	SchemaVersion int `validate:"max=1,min=0"`
	// The EventULID not only makes an event global unique, but also facilitates
	// sorting events into a time ordered sequence.
	EventULID string `validate:"len=26"`
	// The ProxyUserID is our anonymised DrawExact user-key.
	ProxyUserID string `validate:"uuid4"`
	TimeUTC     string `validate:"datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	// The Visit (count) tells us how many cumulative DrawExact sessions the user had launched at the time of this event.
	Visit int `validate:"max=100000,min=1"`
	// The Event tells us what event occurred. E.g. the user completely training step 3.
	Event string `validate:"max=40,min=4"`
	// Some events require additional parameters to adequately describe them.
	// Each event type conveys those parameters in any way it thinks fit encoded into this string.
	Parameters string `validate:"max=80"`
}
```

## The TelemetryEvent data type

```
type TelemetryEvent struct {
  "HowManyPeopleHave": {
    "Launched": 0,
    "LoadedAnExample": 0,
    "TriedToSignIn": 0,
    "SucceededSigningIn": 0,
    "CreatedTheirOwnDrawing": 0,
    "RetreivedTheirASavedDrawing": 0
  },
  "TotalRecoverableErrors": 0,
  "TotalFatalErrors": 0
}
```