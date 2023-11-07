create:
	migrate create -ext sql -dir db/migrations -seq initial-schema

.PHONY: create