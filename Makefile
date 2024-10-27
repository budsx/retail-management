docker-up:
	docker-compose up --build

migrate:
	migrate -path migrations -database "postgres://jubelio:jubeliotest@localhost:5432/retails?sslmode=disable" up

docker-down:
	docker-compose down -v