# Make a DRY storage path synthesis helper

## Context

- Later in the evolution of this App, we will need to store telemetry events and to specify a storage path string that implies time buckets.
- The path format is defined in ./specs/010-init.md

## Objective

- To write the code for a DRY storage path generator derived from a given TelemetryEvent

# Instructions

- Write a function: FormatStoragePath(event *TelemetryEvent) (path, error)
- Write tests for that function