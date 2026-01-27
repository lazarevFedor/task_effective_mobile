build:
	docker compose -f docker/docker-compose.yml --env-file configs/.env up -d --build subscriptions_db
	docker compose -f docker/docker-compose.yml --env-file configs/.env up -d --build server