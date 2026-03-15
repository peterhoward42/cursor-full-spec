# cursor-full-spec

A Google Cloud Function (Go) that provides an HTTP API for telemetry ingestion and analysis reporting. It accepts telemetry events via POST, stores them (e.g. in GCS), and returns an aggregate analysis report via GET.

## What it does

- **POST** — Ingest a single telemetry event (JSON). Payload is validated, serialised as NDJSON gzip, and stored under a path derived from the event’s time and ULID. Responds with `201 Created` on success.
- **GET** — Return the analysis report as JSON: distinct-user counts per event type (e.g. launched, loaded-example, sign-in) and totals for recoverable/fatal JS errors.
- **OPTIONS** — CORS preflight: allows `GET`, `POST`, `OPTIONS` and `Content-Type`, with `Access-Control-Allow-Origin: *`.

Event schema and validation follow the overview-context specification (e.g. `SchemaVersion`, `EventULID`, `ProxyUserID`, `TimeUTC`, `Visit`, `Event`, `Parameters`).

## Build and run

- **Tests:** `go test ./...`
- **Lint:** `golangci-lint run`
- **Deploy:** Deploy as a Gen2 HTTP-triggered function (Go 1.21). Project and region are configured via Make variables.
- **Example requests:** The Makefile provides targets that print ready-to-run `curl` commands for POST (send a sample telemetry event) and GET (fetch the report). Override `BASE_URL` when pointing at your deployed function.

See the [Makefile](Makefile) for the exact commands and variables (`make` or `make help`).

## Layout

- `cmd/function/` — Entry point; wires the app with fake storer/getter for local runs (production would use GCS-backed implementations).
- `internal/app/` — Application logic: HTTP dispatch, validation, serialisation, storage path formatting, event storage/retrieval interfaces, and analysis (distinct-user and error counts).
