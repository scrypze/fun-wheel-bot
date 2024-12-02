# Makefile для управления проектом fun-wheel-bot

# Переменные
DOCKER_COMPOSE=docker compose

# Цели
.PHONY: up down restart update

up:
	@echo "Запуск приложения..."
	$(DOCKER_COMPOSE) up -d --build

down:
	@echo "Остановка приложения..."
	$(DOCKER_COMPOSE) down

restart: down up

update:
	@echo "Обновление репозитория и перезапуск приложения..."
	git pull origin main
	$(MAKE) restart 