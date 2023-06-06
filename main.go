package main

import (
	"database/sql"
	"log"

	"github.com/silaselisha/bank-api/api"
	db "github.com/silaselisha/bank-api/db/sqlc"
	_"github.com/lib/pq"
)

const (
	driver_name = "postgres"
	data_source = "postgresql://root:esilas@localhost:5431/jpmorgan?sslmode=disable"
	address = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(driver_name, data_source)
	if err != nil {
		log.Fatal("Cannot connect to database!")
		return
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("Cannot connect to database!")
		return
	}
}