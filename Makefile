.PHONY: up down restart

up:
	docker compose up -d --build

down:
	docker compose down

restart: down up
	@echo "Контейнеры перезапущены"