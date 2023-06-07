package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/silaselisha/bank-api/api"
	db "github.com/silaselisha/bank-api/db/sqlc"
	"github.com/silaselisha/bank-api/db/utils"
)


func main() {
	var err error
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load ENV!\n", err)
		return
	}

	conn, err := sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		log.Fatal("Cannot connect to database!")
		return
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot connect to database!")
		return
	}
}