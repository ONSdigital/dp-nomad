SHELL=bash

test:
	go test -race -cover ./...
.PHONY: test

audit:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore
.PHONY: audit

build:
	go build ./...
.PHONY: build