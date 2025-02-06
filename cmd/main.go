package main

import (
	"database/sql"
	"log"

	"github.com/flexGURU/simplebank/api"
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/utils"
	_ "github.com/lib/pq"
)



func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("error loading  config: ", err)
	}



	connDb, err := sql.Open(config.DBDriver,config.DSN)
	if err != nil {
		log.Fatal("error opening the database",err)
	}

	store := db.NewStore(connDb)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.StartServer(config.ServerAddress); err != nil {
		log.Fatal("error starting up the server")
	}

}