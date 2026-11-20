SHELL := /bin/sh

API_DIR := elzzes-calculator-api
WEB := elzzes-calculator-web
UI := @elzzes/ui
PNPM ?= pnpm

.DEFAULT_GOAL := help

.PHONY: help install install-api install-web dev dev-api dev-web storybook \
	fmt fmt-check vet go-test go-coverage go-race go-build \
	lint-web typecheck-web test-web coverage-web build-web \
	lint-ui typecheck-ui test-ui coverage-ui build-storybook \
	lint typecheck test coverage build verify \
	compose-up compose-down compose-verify clean

help:
	@printf '%s\n' \
		'Elzzes Calculator - local commands (Docker is optional)' \
		'' \
		'  make install        Install pnpm dependencies and Go modules' \
		'  make dev            Run the API (:8080) and Vite dev server (:5173)' \
		'  make dev-api        Run the Go API only' \
		'  make dev-web        Run the Vite dev server only' \
		'  make storybook      Run Storybook (:6006)' \
		'  make verify         Run every non-Docker check (same as CI)' \
		'  make compose-up     Run the app with Docker Compose (:8080)' \
		'  make compose-verify Build, smoke-test, and stop the Compose stack' \
		'  make compose-down   Stop the Compose stack'

install: install-web install-api

install-web:
	$(PNPM) install --frozen-lockfile

install-api:
	cd $(API_DIR) && go mod download

dev:
	@set -eu; \
	( cd $(API_DIR) && exec go run ./cmd/api ) & api_pid=$$!; \
	trap 'kill $$api_pid 2>/dev/null || true' EXIT INT TERM; \
	$(PNPM) --filter $(WEB) dev

dev-api:
	cd $(API_DIR) && go run ./cmd/api

dev-web:
	$(PNPM) --filter $(WEB) dev

storybook:
	$(PNPM) --filter $(UI) storybook

fmt:
	cd $(API_DIR) && gofmt -w .

fmt-check:
	@cd $(API_DIR) && files="$$(gofmt -l .)"; \
	if [ -n "$$files" ]; then printf '%s\n' "$$files"; echo "gofmt: files need formatting"; exit 1; fi

vet:
	cd $(API_DIR) && go vet ./...

go-test:
	cd $(API_DIR) && go test ./...

go-coverage:
	cd $(API_DIR) && go test ./... -coverprofile=coverage.out
	./scripts/check-go-coverage.sh $(API_DIR)/coverage.out

go-race:
	cd $(API_DIR) && go test -race ./...

go-build:
	cd $(API_DIR) && go build -o /tmp/elzzes-calculator-api ./cmd/api

lint-web:
	$(PNPM) --filter $(WEB) lint

typecheck-web:
	$(PNPM) --filter $(WEB) typecheck

test-web:
	$(PNPM) --filter $(WEB) test

coverage-web:
	$(PNPM) --filter $(WEB) test:coverage

build-web:
	$(PNPM) --filter $(WEB) build

lint-ui:
	$(PNPM) --filter $(UI) lint

typecheck-ui:
	$(PNPM) --filter $(UI) typecheck

test-ui:
	$(PNPM) --filter $(UI) test

coverage-ui:
	$(PNPM) --filter $(UI) test:coverage

build-storybook:
	$(PNPM) --filter $(UI) build-storybook

lint: fmt-check vet lint-web lint-ui

typecheck: typecheck-web typecheck-ui

test: go-test test-web test-ui

coverage: go-coverage coverage-web coverage-ui

build: go-build build-web build-storybook

verify: fmt-check vet go-test go-coverage go-race \
	lint-web typecheck-web test-web coverage-web \
	lint-ui typecheck-ui test-ui coverage-ui \
	build-web build-storybook go-build

compose-up:
	docker compose up --build

compose-down:
	docker compose down

compose-verify:
	@set -eu; \
	trap 'docker compose down' EXIT; \
	docker compose up --build -d; \
	curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://localhost:8080/healthz; \
	echo; \
	curl -fsS http://localhost:8080/ > /tmp/elzzes-calculator-index.html; \
	head -n 5 /tmp/elzzes-calculator-index.html; \
	curl -sS http://localhost:8080/api/v1/calculate -H 'Content-Type: application/json' \
		-d '{"operation":"divide","operands":[10,0]}'; \
	echo

clean:
	rm -rf $(API_DIR)/coverage.out $(WEB)/coverage $(WEB)/dist $(UI)/coverage $(UI)/storybook-static
