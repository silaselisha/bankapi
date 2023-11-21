package main

import (
	"database/sql"
	"log"

	"github.com/silaselisha/bankapi/api"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	envs, err := utils.Load(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(envs.DBdriver, envs.DBsource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	server.Start(envs.Address)
}
