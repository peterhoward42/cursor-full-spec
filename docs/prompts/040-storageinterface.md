# Introduce a storage interface dependency

## Context

- The existing code already supports DI of external system interfaces using the Dependencies struct.
- This task will introduce an EventStorer interface and add it to Dependencies

## Objective

- To define an interface for storing TelemetryEvents
- To introduce the interface to the Dependency struct
- To implement a fake, test double implementation of the interface

## Non objectives
- Do not implement a real world / production implementation of the interface yet

# Instructions

- Define the following "EventStorer" interface:

```
StoreEventIfNotExists(event *TelemetryEvent) error
```

- Generate the code for a fake test-double implementation of the EventStorer interface
- Update the Dependencies structure to have a EventStorer field
- Use the fake test double implemenation in all the existing Dependencies construction sites