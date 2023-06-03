// ** Entry point to every test
// ** Mocking a database
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_"github.com/lib/pq"
)

const (
	driver_name = "postgres"
	data_source = "postgres://root:esilas@localhost:5431/jpmorgan?sslmode=disable"
)

var testQueries *Queries
var conn *sql.DB

func TestMain(m *testing.M) {
	var err error
	conn, err = sql.Open(driver_name, data_source)

	if err != nil {
		log.Fatal("database connection unsuccessfully!")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
