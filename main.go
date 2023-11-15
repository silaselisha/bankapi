package main

import (
	"database/sql"
	"log"

	"github.com/silaselisha/bankapi/api"
	db "github.com/silaselisha/bankapi/database/sqlc"

	_ "github.com/lib/pq"
)

const (
	ADDRESS   = ":8080"
	DB_SOURCE = "postgresql://root:esilas@localhost:5432/bankapi?sslmode=disable"
)

func main() {
	conn, err := sql.Open("postgres", DB_SOURCE)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	server.Start(ADDRESS)
}
