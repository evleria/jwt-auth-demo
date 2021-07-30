swag:
	swag init --parseDependency --parseDepth=5
compose-build:
	docker-compose build
compose-up:
	docker-compose up -d
compose-down:
	docker-compose down

.PHONY: swag, compose-build, compose-up, compose-down