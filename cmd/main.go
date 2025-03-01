package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net"

	"github.com/flexGURU/simplebank/api"
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/gapi"
	"github.com/flexGURU/simplebank/mail"
	"github.com/flexGURU/simplebank/pb"
	"github.com/flexGURU/simplebank/utils"
	"github.com/flexGURU/simplebank/worker"
	"github.com/hibiken/asynq"
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

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	redisDistro := worker.NewRedisTaskDistributer(redisOpt)

	go runTaskProcessor(redisOpt, store, config)
	startGinServer(config, store, redisDistro)

}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config utils.Config)  {

	mailSender := mail.NewGmailSender(config.EmailSendName, config.From_Email, config.EamilPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailSender)

	slog.Info(
		"task taskProcessor started",
	)
	if err := taskProcessor.Start(); err != nil {
		slog.Error("failed to start task processor %w", err)
	}
	
	
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

func startGinServer(config utils.Config,store db.Store, taskDistributer worker.TaskDistributer)  {
	server, err := api.NewServer(config, store, taskDistributer)
	if err != nil {
		log.Fatal(err)
	}
	err = server.StartServer(config.HTTPServerAddress); 
	slog.Info("started the server", 
				slog.String("port", config.HTTPServerAddress),
				)
	if err != nil {
		log.Fatal("error starting up the server")
	}
}