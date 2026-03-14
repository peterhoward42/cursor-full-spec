# Make a DRY TelemetryEvent JSON payload validator

## Overview

- The app contains a request handler that will receive a single JSON encoded TelemetryEvent.
- The request handler will need to parse and validate that event and report errors.


# Instructions

- Write encapsulated code that receives a string and attempts to parse it into a         TelemetryEvent, with comprehensive validation.
- Write tests for that code
- Do not wire it in to the rest of the application yet