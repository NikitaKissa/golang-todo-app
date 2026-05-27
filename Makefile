include .env
export

PROJECT_ROOT := $(CURDIR)
export PROJECT_ROOT

SHELL := bash

env-up:
	docker compose up -d todoapp-postgres

env-down:
	docker compose down todoapp-postgres

env-port-forward:
	docker compose up -d port-forwarder

env-port-close:
	docker compose down port-forwarder

migrate-create:
	@powershell -Command "if (-not '$(seq)') { Write-Host 'seq is empty'; exit 1 }" && \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	docker compose run --rm todoapp-postgres-migrate \
    		-path /migrations \
    		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
    		$(action)

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down