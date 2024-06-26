GO_SRCS := $(shell find . -type f -name '*.go' -a ! \( -name 'zz_generated*' -o -name '*_test.go' \))
GO_TESTS := $(shell find . -type f -name '*_test.go')
TAG_NAME = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
TAG_NAME_DEV = $(shell git describe --tags --abbrev=0 2>/dev/null)
VERSION_CORE = $(shell echo $(TAG_NAME)')
VERSION_CORE_DEV = $(shell echo $(TAG_NAME_DEV)')
GIT_COMMIT = $(shell git rev-parse --short=7 HEAD)
VERSION = $(or $(and $(TAG_NAME),$(VERSION_CORE)),$(and $(TAG_NAME_DEV),$(VERSION_CORE_DEV)-dev),$(GIT_COMMIT))
DB_DSN ?= $(shell cat .env | grep DB_DSN | cut -d '=' -f 2)
DB_DSN := $(or $(DB_DSN),$(shell cat .env.local | grep DB_DSN | cut -d '=' -f 2))
E2E_TESTS := $(shell find e2e -type f -name '*.tape')

MIGRATIONS := $(shell find ./migrations -type f -name '*.sql')
MIGRATION_NAME?=migration_name

golint := $(shell which golangci-lint)
ifeq ($(golint),)
golint := $(shell go env GOPATH)/bin/golangci-lint
endif

pkgsite := $(shell which pkgsite)
ifeq ($(pkgsite),)
pkgsite := $(shell go env GOPATH)/bin/pkgsite
endif

mockery := $(shell which mockery)
ifeq ($(mockery),)
mockery := $(shell go env GOPATH)/bin/mockery
endif

vhs := $(shell which vhs)
ifeq ($(vhs),)
vhs := $(shell go env GOPATH)/bin/vhs
endif

goose := $(shell which goose)
ifeq ($(goose),)
goose := $(shell go env GOPATH)/bin/goose
endif

sqlc := $(shell which sqlc)
ifeq ($(sqlc),)
sqlc := $(shell go env GOPATH)/bin/sqlc
endif

bin/toshokan: $(GO_SRCS)
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/checksums.txt: $(addprefix bin/,$(bins))
	sha256sum -b $(addprefix bin/,$(bins)) | sed 's/bin\///' > $@

bin/checksums.md: bin/checksums.txt
	@echo "### SHA256 Checksums" > $@
	@echo >> $@
	@echo "\`\`\`" >> $@
	@cat $< >> $@
	@echo "\`\`\`" >> $@

bin/toshokan-darwin-amd64: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/toshokan-darwin-arm64: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/toshokan-freebsd-amd64: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/toshokan-freebsd-arm64: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/toshokan-linux-amd64: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/toshokan-linux-arm64: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/toshokan-linux-riscv64: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bin/toshokan-windows-amd64.exe: $(GO_SRCS)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

bins := toshokan-darwin-amd64 toshokan-darwin-arm64 toshokan-freebsd-arm64 toshokan-freebsd-arm64 toshokan-linux-amd64 toshokan-linux-arm64 toshokan-linux-riscv64 toshokan-windows-amd64.exe

.PHONY: build-all
build-all: $(addprefix bin/,$(bins)) bin/checksums.md

.PHONY: unit
unit:
	go test -race -covermode=atomic -tags=unit -timeout=30s ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: migration
migration:
	$(goose) -dir db/migrations -s create $(MIGRATION_NAME) sql

.PHONY: up
up: $(MIGRATIONS)
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	@$(goose) -dir db/migrations postgres $(DB_DSN) up

.PHONY: drop
drop:
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	@$(goose) -dir db/migrations postgres $(DB_DSN) reset

.PHONY: sql
sql: $(sqlc)
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	@DB_DSN=$(DB_DSN) $(sqlc) generate

$(golint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(pkgsite):
	go install golang.org/x/pkgsite/cmd/pkgsite@latest

$(vhs):
	go install github.com/charmbracelet/vhs@latest

$(goose):
	go install -tags 'no_clickhouse,no_mssql,no_mysql,no_turso,no_vertica,no_ydb' github.com/pressly/goose/v3/cmd/goose

$(sqlc):
	go install github.com/sqlc-dev/sqlc/cmd/sqlc

.PHONY: license
license:
	go run ./licensing.go


.PHONY: lint
lint: $(golint)
	$(golint) run ./...

%.result: %.tape
	@echo "Running e2e test for $<..."
	mkdir -p e2e/tmp
	cd e2e/tmp && PATH=$(shell pwd)/bin:$(PATH) $(vhs) ../../$<
	echo "pass" > $@

.PHONY: e2e
e2e: bin/toshokan $(vhs) $(E2E_TESTS:.tape=.result)

.PHONY: mocks
mocks:
	$(mockery)

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: version
version:
	@echo VERSION_CORE=${VERSION_CORE}
	@echo VERSION_CORE_DEV=${VERSION_CORE_DEV}
	@echo VERSION=${VERSION}

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

.PHONY: doc
doc: $(pkgsite)
	$(pkgsite) -open .
