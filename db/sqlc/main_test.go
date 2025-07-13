package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/volskyi-dmytro/st-bank/util"
)

// testQueries is a global variable to hold database queries for testing
var testQueries *Queries

// testDB is a global variable to hold database connection for testing
var testDB *sql.DB

// TestMain sets up the test database connection and runs all tests
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can not connect to the database: ", err)
	}

	testDB = conn
	testQueries = New(conn)

	os.Exit(m.Run())
}
