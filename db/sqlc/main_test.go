// ** Entry point to every test
// ** Mocking a database
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/silaselisha/bank-api/db/utils"
)

var testQueries *Queries
var conn *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..")

	if err != nil {
		log.Fatal("Could not load ENVs!\n", err)
	}
	
	conn, err = sql.Open(config.DBdriver, config.DBsource)

	if err != nil {
		log.Fatal("database connection unsuccessfully!")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
