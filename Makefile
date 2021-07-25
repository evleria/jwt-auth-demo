swag:
	swag init --parseDependency --parseDepth=5 -d ./cmd/rest -o ./cmd/rest/docs
compose-build:
	docker-compose build
compose-up:
	docker-compose up -d
compose-down:
	docker-compose down

.PHONY: swag, compose-build, compose-up, compose-down