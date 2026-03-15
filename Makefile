# Makefile for cursor-full-spec
# Run "make" or "make help" to see available targets.

.PHONY: test lint deploy post get

test:
	go test ./...

lint:
	golangci-lint run

# Deploy the Google Cloud Function using source deployment.
# Set PROJECT_ID and REGION (e.g. make deploy PROJECT_ID=my-project REGION=europe-west1).
deploy:
	gcloud functions deploy Function \
		--gen2 \
		--runtime=go121 \
		--region=$${REGION:-REGION} \
		--project=$${PROJECT_ID:-PROJECT_ID} \
		--source=. \
		--entry-point=Function \
		--trigger-http

# Show a command line that POSTs a valid telemetry event to the API.
# Set BASE_URL when running (e.g. BASE_URL=https://... make post) to print a ready-to-run curl.
post:
	@echo "POST a valid telemetry event:"
	@echo "curl -X POST \"$${BASE_URL:-https://REGION-PROJECT_ID.cloudfunctions.net/Function}\" \\"
	@echo "  -H 'Content-Type: application/json' \\"
	@echo "  -d '{\"SchemaVersion\":1,\"EventULID\":\"01ARZ3NDEKTSV4RRFFQ69G5FAV\",\"ProxyUserID\":\"550e8400-e29b-41d4-a716-446655440000\",\"TimeUTC\":\"2025-03-15T12:00:00Z\",\"Visit\":1,\"Event\":\"launched\",\"Parameters\":\"\"}'"

# Show a command line that GETs the analysis report from the API.
# Set BASE_URL when running (e.g. BASE_URL=https://... make get) to print a ready-to-run curl.
get:
	@echo "GET the analysis report:"
	@echo "curl \"$${BASE_URL:-https://REGION-PROJECT_ID.cloudfunctions.net/Function}\""
