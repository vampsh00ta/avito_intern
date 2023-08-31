test:
	go test ./internal/transport
dev:
	air

run:
	docker-compose build && docker-compose up
