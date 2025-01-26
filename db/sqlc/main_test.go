package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_"github.com/lib/pq"
)

var testQueries *Queries
const (
	dbDriver = "postgres"
	dsn = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
)

var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error

	dbConn, err = sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal("problem opening db", err)
	}

	defer dbConn.Close()

	testQueries = New(dbConn)

	os.Exit(m.Run())


	





}