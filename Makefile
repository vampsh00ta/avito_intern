test:
	go test ./internal/transport
dev:
	air
prod:
	docker-compose up -d

run:
	docker-compose build && docker-compose up
