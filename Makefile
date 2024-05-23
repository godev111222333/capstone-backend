gen:
	go generate ./...

dev-up:
	@docker compose \
		-f docker/docker-compose.yml up -d

dev-down:
	@docker compose \
		-f docker/docker-compose.yml down

migrate-up:
	@migrate -path migration -database "postgres://root:password@localhost:5432/rental_car?sslmode=disable" -verbose up

migrate-down:
	@migrate -path migration -database "postgres://root:password@localhost:5432/rental_car?sslmode=disable" -verbose down

dev-run-core:
	@go run src/cmd/core/main.go

dev-run-api:
	@go run src/cmd/api/main.go
