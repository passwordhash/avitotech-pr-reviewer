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
