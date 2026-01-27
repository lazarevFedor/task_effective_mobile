.PHONY: build

build:
	docker compose -f deployments/docker-compose.yml --env-file configs/.env up -d --build subscriptions_db
	docker compose -f deployments/docker-compose.yml --env-file configs/.env up -d --build server