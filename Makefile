gen:
	go generate ./...
	mockgen -source src/service/push_notification_service.go -destination=src/service/push_notification_service_mock.go -package=service

import:
	go run src/cmd/import/main.go

test:
	@go test -v ./...

dev-build:
	@docker compose build

dev-up:
	@docker compose \
		-f docker-compose.yml up -d

dev-down:
	@docker compose \
		-f docker-compose.yml down

migrate-up:
	@migrate -path migration -database "postgres://root:password@localhost:5432/rental_car?sslmode=disable" -verbose up
	go run src/cmd/import/main.go

migrate-down:
	@migrate -path migration -database "postgres://root:password@localhost:5432/rental_car?sslmode=disable" -verbose down

dev-run:
	@go run src/cmd/api/main.go
