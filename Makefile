swag:
	swag init -g cmd/main.go

database:
	sudo docker run --name=walletter-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres

migrate.up:
	migrate -path ./schema -database postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable up

migrate.down:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' down

run:
	sudo docker-compose up