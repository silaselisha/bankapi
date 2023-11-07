package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_"github.com/lib/pq"
)

var testQueries *Queries
func TestMain(m *testing.M) {
	db, err := sql.Open("postgres", "postgres://root:esilas@localhost:5432/bankapi?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(db)
	os.Exit(m.Run())
}
