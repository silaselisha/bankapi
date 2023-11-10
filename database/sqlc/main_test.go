package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	DB_DRIVER = "postgres"
	DB_SOURCE = "postgresql://root:esilas@localhost:5432/bankapi?sslmode=disable"
)

var testQueries *Queries
var conn *sql.DB

func TestMain(m *testing.M) {
	var err error
	conn, err = sql.Open(DB_DRIVER, DB_SOURCE)
	// fmt.Printf("type: %T -> value: %v\n", conn, conn)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
