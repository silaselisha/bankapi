create:
	migrate create -ext sql -dir db/migrations -seq indexing

generate:
	sqlc generate

postgres:
	docker run --name postgres16 -dp 5432:5432 -e POSTGRES_PASSWORD=esilas -e POSTGRES_DB=bankapi -e POSTGRES_USER=root postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root bankapi

dropdb:
	docker exec -it postgres16 dropdb bankapi

migrateup:
	migrate -path db/migrations -database "postgres://root:success87@postgres.c0yazamau9ss.us-east-1.rds.amazonaws.com:5432/bankapi" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://root:success87@postgres.c0yazamau9ss.us-east-1.rds.amazonaws.com:5432/bankapi" -verbose down

gotest:
	go test -v -cover ./...

.PHONY: create generate postgres createdb dropdb migrateup migratedown gotest