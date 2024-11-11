GO_SRCS := $(shell find . -type f -name '*.go' -a -name '*.tpl' -a ! \( -name 'zz_generated*' -o -name '*_test.go' \))
GO_TESTS := $(shell find . -type f -name '*_test.go')
TAG_NAME = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
TAG_NAME_DEV = $(shell git describe --tags --abbrev=0 --match "scan/plugin/anchira/*" 2>/dev/null)
VERSION_CORE = $(shell echo $(TAG_NAME) | sed 's|scan/plugin/anchira/\(.*\)|\1|')
VERSION_CORE_DEV = $(shell echo $(TAG_NAME_DEV) | sed 's|scan/plugin/anchira/\(.*\)|\1|')
GIT_COMMIT = $(shell git rev-parse --short=7 HEAD)
VERSION = $(or $(and $(TAG_NAME),$(VERSION_CORE)),$(and $(TAG_NAME_DEV),$(VERSION_CORE_DEV)-dev),$(GIT_COMMIT))

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

golint := $(shell which golangci-lint)
ifeq ($(golint),)
golint := $(shell go env GOPATH)/bin/golangci-lint
endif

goreleaser := $(shell which goreleaser)
ifeq ($(goreleaser),)
goreleaser := $(shell go env GOPATH)/bin/goreleaser
endif

build: dist/tsk_anchira_$(VERSION)_$(GOOS)_$(GOARCH).so

.PHONY: dist/tsk_anchira_$(VERSION)_linux_amd64.so
dist/tsk_anchira_$(VERSION)_linux_amd64.so: $(GO_SRCS)
	GOARCH=amd64 \
	CGO_ENABLED=1 go build -buildmode=plugin -trimpath -ldflags '-X main.version=${VERSION} -s -w' -o "$@" ./main.go

.PHONY: dist/tsk_anchira_$(VERSION)_linux_arm64.so
dist/tsk_anchira_$(VERSION)_linux_arm64.so: $(GO_SRCS)
	GOARCH=arm64 \
	CGO_ENABLED=1 go build -buildmode=plugin -trimpath -ldflags '-X main.version=${VERSION} -s -w' -o "$@" ./main.go

.PHONY: snapshot
snapshot: $(goreleaser)
	$(goreleaser) release --snapshot --clean

.PHONY: release
release: $(goreleaser)
	$(goreleaser) release --clean

.PHONY: unit
unit:
	go test -race -covermode=atomic -tags=unit -timeout=30s ./...

.PHONY: generate
generate:
	go generate ./...

$(golint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(goreleaser):
	go install github.com/goreleaser/goreleaser/v2@latest

.PHONY: lint
lint: $(golint)
	$(golint) run ./...

.PHONY: clean
clean:
	rm -rf dist/
