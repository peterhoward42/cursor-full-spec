# Put in a Makefile

## Rational

Given a codebase a Makefile has the following benefits

- Quick and easy process automation for the developer
- Provides how-to documentation as a useful side effect
- A customary way for consumers of the repo to see quickly how to shape illustrative http requests from the command line

## Instructions

Create a Makefile with targets as follows:

- A "test" target that simply runs `go test ./...`
- A "lint" target that runs `golangci-lint`
- A "deploy" target that fires the command line command to deploy the google cloud function - using the "source" deployment method
    - If you don't have sufficient config data to generate this code accurately - use illustrative placeholder elements in the generated command line
- A "post" target that shows a command line that POSTS a valid telemetry event to the API
- A "get" target that shows a command line that makes a GET request to the API to receive a report.

