.PHONY: lint

lint:
	golint `go list ./...`
