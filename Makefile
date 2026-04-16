.PHONY: fmt lint check

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

check: fmt lint
