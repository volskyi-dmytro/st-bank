package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// Database connection constants
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/st_bank?sslmode=disable"
)

// testQueries is a global variable to hold database queries for testing
var testQueries *Queries

// testDB is a global variable to hold database connection for testing
var testDB *sql.DB

// TestMain sets up the test database connection and runs all tests
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can not connect to the database: ", err)
	}

	testDB = conn
	testQueries = New(conn)

	os.Exit(m.Run())
}
