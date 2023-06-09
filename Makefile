docker:
	docker run --name postgres15 -dp 5431:5432 -e POSTGRES_PASSWORD=esilas -e POSTGRES_USER=root postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root jpmorgan

dropdb:
	docker exec -it postgres15 dropdb --username=root jpmorgan

migrate:
	migrate create -ext sql -dir db/migration -seq schema 

migrateup: 
	migrate -path db/migration -database "postgres://root:esilas@localhost:5431/jpmorgan?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:esilas@localhost:5431/jpmorgan?sslmode=disable" -verbose down

migrateupone: 
	migrate -path db/migration -database "postgres://root:esilas@localhost:5431/jpmorgan?sslmode=disable" -verbose up 1

migratedownone:
	migrate -path db/migration -database "postgres://root:esilas@localhost:5431/jpmorgan?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v ./...

server:
	go run main.go

.PHONY: docker createdb dropdb migrate migrateup migratedown sqlc server

