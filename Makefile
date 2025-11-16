ENV_FILE ?= .env.local

.EXPORT_ALL_VARIABLES:
include $(ENV_FILE)

# ==========================
# Локальная разработка / Тесты
# ==========================
up:
	docker compose \
		--env-file .env.local \
		-p pr-reviewer \
		up -d --build

# ==========================
# Миграции
# ==========================

migration-create:
	# Пример использования: make migration-create name=create_users_table
	@if [ -z "$(name)" ]; then \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(name)

migration-up:
	migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose up

migration-down:
	migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose down

# ==========================
# e2e тесты
# ==========================
HTTP_PORT_E2E ?= 8081
CONTAINER_NAME_E2E ?= pr-reviewer-e2e

test-e2e-up:
	docker compose \
		-f docker-compose.e2e.yml \
		-p $(CONTAINER_NAME_E2E) \
		up -d --build

test-e2e: test-e2e-up
	go mod tidy
	HTTP_PORT=$(HTTP_PORT_E2E) go test -v ./tests/...
	$(MAKE) test-e2e-down

test-e2e-down:
	docker compose \
		-f docker-compose.e2e.yml \
		-p $(CONTAINER_NAME_E2E) \
		down -v
