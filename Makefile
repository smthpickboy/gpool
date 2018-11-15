GO=go

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: check

.PHONY: all check lint test

lint:
	@gometalinter --config=.gometalint ./...

PACKAGES = $(shell go list ./...|grep -v /vendor/)
test: check
	$(GO) test ${PACKAGES}

cov: check
	gocov test $(PACKAGES) | gocov-html > coverage.html

check: lint
	@$(GO) tool vet ${SRC}

hook:
	find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
