create:
	migrate create -ext sql -dir db/migrations -seq initial-schema

generate:
	sqlc generate

migrateup:
	migrate -path db/migrations -database "postgres://postgres:esilas@localhost:5432/bankapi?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://postgres:esilas@localhst:5432/bankapi?sslmode=disable" -verbose down

.PHONY: create generate migrateup