test:
	go test ./internal/transport
hello:
	echo "Hello"
dev:
	air
migrate:
	 migrate -path migrations -database 'postgres://avito:avito@db:5433/avito?sslmode=disable' up

docker-up:
	docker-compose up
docker-build:
	docker-compose build
