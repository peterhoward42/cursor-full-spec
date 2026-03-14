# Introduce an EventGetter interface dependency

## Context

- The existing code already supports DI of external system interfaces using the Dependencies struct.
- This task will introduce an EventGetter interface and add it to Dependencies

## Objective

- To define an interface for getting all stored TelemetryEvents
- To add an EventGetter field to the Dependency struct
- To implement a fake, test double implementation of the interface

## Non objectives
- Do not implement a real world / production implementation of the interface yet

# Instructions

- Define the following "EventGetter" interface:

```
GetAllStoredEvents() (event *[]TelemetryEvent), error)
```

- Generate the code for a fake test-double implementation of the EventGetter interface
- Update the Dependencies structure to have an EventGetter field
- Use the fake test double implemenation in all the existing Dependencies construction sites