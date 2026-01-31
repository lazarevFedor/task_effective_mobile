.PHONY: build

build:
	docker compose -f deployments/docker-compose.yml --env-file configs/.env up -d --build subscriptions_db
	docker compose -f deployments/docker-compose.yml --env-file configs/.env up -d --build server

stop:
	docker compose -f deployments/docker-compose.yml --env-file configs/.env down subscriptions_db
	docker compose -f deployments/docker-compose.yml --env-file configs/.env down server

up:
	migrate -database 'postgres://user:1234@localhost:5432/subscriptions_db?sslmode=disable' -path internal/migrations up

down:
	migrate -database 'postgres://user:1234@localhost:5432/subscriptions_db?sslmode=disable' -path internal/migrations down

swagger:
	swag init -g internal/server/server.go -o api