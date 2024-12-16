.PHONY: run down restart

run:
	docker compose up -d --build

down:
	docker compose down

restart: down run
	@echo "Контейнеры перезапущены"