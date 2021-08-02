swag:
	swag init --parseDependency --parseDepth=5
compose-build:
	docker-compose build
compose-up:
	docker-compose up -d
compose-down:
	docker-compose down
lint:
	golangci-lint run ./internal/... && golangci-lint run ./main.go
import:
	goimports -local "github.com/evleria/jwt-auth-demo" -w .

.PHONY: swag, compose-build, compose-up, compose-down, lint, import