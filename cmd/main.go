package main

import (
	"database/sql"
	"log"

	"github.com/flexGURU/simplebank/api"
	db "github.com/flexGURU/simplebank/db/sqlc"
	_"github.com/lib/pq"

)

const (
	dbDriver = "postgres"
	dsn      = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
	address = ":8080"
)

func main() {

	connDb, err := sql.Open(dbDriver,dsn)
	if err != nil {
		log.Fatal("error opening the database")
	}

	store := db.NewStore(connDb)

	server := api.NewServer(store)
	if server.StartServer(address); err != nil {
		log.Fatal("error starting up the server")
	}

}