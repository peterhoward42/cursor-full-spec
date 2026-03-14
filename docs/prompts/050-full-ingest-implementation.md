# Replacing the placeholder Application.IngestEvent function behaviour with
  real behaviour

## Objectives

- To implement the following behaviour in the Appliction.IngestEvent
    - Parse and validate the request payload (use the existing ParseAndValidateTelemetryEvent)
    - Serialise the validated TelemetryEvent as NDJSON gzip compressed
    - Write the serialised event to storage using the existing GCSEventStorer
    - Use the existing FormatStoragePath function to generate the storage path argument required by the EventStorer

## Instructions

- Generate the code as indicated by the given Objectives above
- Write tests for Application.IngestEvent with coverage defined as follows
    -  Validate that the sequence of operations defined above has been properly orchestrated and that the correct data has propagated through
    the call stack
    - Validate all the new branching logic introduced by this change
    - Do NOT write tests that duplicate fully the validation performed by unit tests elsewhere - for example the tests for storage path validation