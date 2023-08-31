.PHONY: test up down migrate-new deps

up-docker:
	docker-compose up -d --build --force-recreate

down-docker:
	docker-compose down

test:
	go test -v -race ./...

up: down
	docker run --name astrologer-db-test -p 11453:5432 -e POSTGRES_PASSWORD=dev -d postgres:15.3-alpine

down:
	docker rm -f astrologer-db-test

migrate-new:
	goose -dir ./migrations create $(name) sql

deps:
	go mod download
	go install github.com/pressly/goose/v3/cmd/goose@latest
