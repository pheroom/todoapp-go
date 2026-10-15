include .env
export

PROJECT_ROOT := $(CURDIR)
#PROJECT_ROOT := $(shell pwd)

test:
	@echo $(PROJECT_ROOT)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Clean all todoapp-postgres data? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  	docker compose down todoapp-postgres && \
	  	rm -rf out/pgdata && \
	  	echo "Env files cleaned"; \
	else \
		echo "Cleaning process was canceled"; \
	fi

env-port-forwarder:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "Отсутствует 'seq'"; \
		exit 1; \
	fi; \
	docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(act)" ]; then \
  		echo "Отсутствует 'act'"; \
  		exit 1; \
	fi; \
	docker compose run --rm migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(act)"

migrate-up:
	@make migrate-action act=up

todoapp-run:
	@go run cmd/todoapp/main.go