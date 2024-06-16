.PHONY: up down migrate migrate-up migrate-down migrate-deps migration lint lint-deps imports imports-deps format format-deps

LOCAL_BIN := $(CURDIR)/bin
MIGRATE := $(LOCAL_BIN)/migrate
MIGRATE_VERSION ?= v4.15.2
DB_URL := "mysql://user:user_password@tcp(localhost:3306)/my_database"
MIGRATION_NAME ?= migration

up:
	@docker-compose build
	@docker-compose up -d
	$(MAKE) migrate-up


down:
	@docker-compose down

migrate-deps:
ifeq ($(wildcard $(MIGRATE)),)
	@echo "Installing migrate tool..."
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install -tags 'mysql' -mod=mod github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
endif

migrate-up: migrate-deps
	@echo "Running migrations up..."
	@$(MIGRATE) -path ./migrations -database $(DB_URL) up

migrate-down: migrate-deps
	@echo "Reverting migrations..."
	@$(MIGRATE) -path ./migrations -database $(DB_URL) down

migrate: migrate-up

migration: migrate-deps
	$(MIGRATE) create -dir ./migrations -ext sql $(MIGRATION_NAME)

lint-deps:
ifeq ($(wildcard $(LOCAL_BIN)/golangci-lint),)
	@echo "Installing golangci-lint..."
	@mkdir -p $(LOCAL_BIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) v1.42.1
endif

lint: lint-deps
	@echo "Running linter"
	$(LOCAL_BIN)/golangci-lint run --config .golangci.yml --timeout=3m

imports-deps:
ifeq ($(wildcard $(LOCAL_BIN)/goimports-reviser),)
	@echo "Installing goimports-reviser..."
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/incu6us/goimports-reviser/v3@latest
endif

imports: imports-deps
	@echo "Running imports"
	find . -name \*.go \
	    -exec $(LOCAL_BIN)/goimports-reviser -rm-unused -set-alias \
	    -format ./... \;

format-deps:
ifeq ($(wildcard $(LOCAL_BIN)/gofumpt),)
	@echo "Installing gofumpt..."
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install mvdan.cc/gofumpt@latest
endif

format: format-deps
	@echo "Running gofumpt"
	$(LOCAL_BIN)/gofumpt -w .