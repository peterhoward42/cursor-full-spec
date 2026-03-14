# Introduce a Google Cloud Storage implementation of the EventStorer interface.

## Context

- The existing code already models an EventStorer interface, and contains a fake, test-double implementation.

## Objective
- This tasks generates the code for a Google Cloud Storage implementation of the EventStorer interface

## Non objectives
- Do not wire the new implementation into the App yet.

# Instructions

- Generate the code to implement the EventStorer interface using Google Cloud storage.
- Use the google.com/go/storage Go package
- Do not wire up this implementation yet
