package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/flexGURU/simplebank/api"
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/gapi"
	"github.com/flexGURU/simplebank/pb"
	"github.com/flexGURU/simplebank/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	startGRPCServer(config, store)

	

}
func startGRPCServer(config utils.Config,store db.Store)  {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)


	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("startinng grpc Listener")
	for service := range grpcServer.GetServiceInfo() {
		log.Println("Registered Service:", service)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}

	

}

func startGinServer(config utils.Config,store db.Store)  {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.StartServer(config.HTTPServerAddress); err != nil {
		log.Fatal("error starting up the server")
	}
}