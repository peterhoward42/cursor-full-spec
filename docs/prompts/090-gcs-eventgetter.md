# Introduce a Google Cloud Storage implementation of the EventGetter interface.

## Context

- The existing code already models an EventGetter interface, and contains a fake, test-double implementation.

## Objective
- This tasks generates the code for a Google Cloud Storage implementation of the EventGetter interface

## Non objectives
- Do not wire the new implementation into the App yet.

# Instructions

- Generate the code to implement the EventGetter interface using Google Cloud storage.
- Use the google.com/go/storage Go package
- Do not wire up this implementation yet
