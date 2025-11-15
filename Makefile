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
	POSTGRES_HOST=localhost migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose up

migration-down:
	POSTGRES_HOST=localhost migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose down
