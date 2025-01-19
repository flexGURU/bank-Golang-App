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


func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal("problem opening db", err)
	}

	defer conn.Close()

	testQueries = New(conn)

	os.Exit(m.Run())


	





}