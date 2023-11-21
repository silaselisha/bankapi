package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/silaselisha/bankapi/db/utils"
)

var testQueries *Queries
var conn *sql.DB

func TestMain(m *testing.M) {
	var err error
	envs, err := utils.Load("../..")
	if err != nil {
		log.Fatal(err)
	}
	conn, err = sql.Open(envs.DBdriver, envs.DBsource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}