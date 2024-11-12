DB_DSN ?= $(shell cat .env.local | grep DB_DSN | cut -d '=' -f 2)
DB_DSN := $(or $(DB_DSN),$(shell cat .env | grep DB_DSN | cut -d '=' -f 2))

MIGRATIONS := $(shell find ./migrations -type f -name '*.sql')
MIGRATION_NAME?=migration_name

golint := $(shell which golangci-lint)
ifeq ($(golint),)
golint := $(shell go env GOPATH)/bin/golangci-lint
endif

goose := $(shell which goose)
ifeq ($(goose),)
goose := $(shell go env GOPATH)/bin/goose
endif

sqlc := $(shell which sqlc)
ifeq ($(sqlc),)
sqlc := $(shell go env GOPATH)/bin/sqlc
endif

goreleaser := $(shell which goreleaser)
ifeq ($(goreleaser),)
goreleaser := $(shell go env GOPATH)/bin/goreleaser
endif

.PHONY: build
build:
	$(goreleaser) build --single-target --snapshot --clean

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

.PHONY: migration
migration:
	$(goose) -dir db/migrations -s create $(MIGRATION_NAME) sql

.PHONY: up
up: $(MIGRATIONS)
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	@$(goose) -dir db/migrations postgres "$(DB_DSN)" up

.PHONY: drop
drop:
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	@$(goose) -dir db/migrations postgres "$(DB_DSN)" reset

.PHONY: sql
sql: $(sqlc)
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	@DB_DSN="$(DB_DSN)" $(sqlc) generate

$(golint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(goose):
	go install -tags 'no_clickhouse,no_mssql,no_mysql,no_turso,no_vertica,no_ydb' github.com/pressly/goose/v3/cmd/goose

$(sqlc):
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

$(goreleaser):
	go install github.com/goreleaser/goreleaser/v2@latest

.PHONY: license
license:
	go run ./licensing.go

.PHONY: lint
lint: $(golint)
	$(golint) run ./...

.PHONY: clean
clean:
	rm -rf dist/
