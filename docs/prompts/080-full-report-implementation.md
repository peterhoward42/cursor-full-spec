# Replacing the placeholder Application.AnalysisReport function behaviour with
  real behaviour

## Objectives

- To replace the placeholder code in the Appliction.AnalysisReport with the real behaviour
    - Fetch all the TelemetryEvents from the Dependency-Injected EventGetter
    - Create the analysis to reply with using the Analyser

## Instructions

- Fetch all the TelemetryEvents from the Dependency-Injected EventGetter
- Create the analysis to reply with using the Analyser

- Write tests for Application.AnalysisReport with coverage defined as follows
    -  Validate that the sequence of operations defined above has been properly orchestrated and that the correct data has propagated through
    the call stack
    - Validate all the new branching logic introduced by this change
    - Do NOT write tests that duplicate fully the validation performed by unit tests elsewhere - for example the tests for storage path validation