ifneq ("$(wildcard .env)", "")
	include .env
	export $(shell sed 's/=.*//' .env)
endif


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
	migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose up

migration-down:
	migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)" -verbose down
