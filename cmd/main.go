package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}


func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("error loading  config: ", err)  
	}

	connDb, err := sql.Open(config.DBDriver,config.DSN)
	if err != nil {
		log.Fatal("error opening the database",err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	store := db.NewStore(connDb)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	redisDistro := worker.NewRedisTaskDistributer(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, redisOpt, store, config)
	startGinServer(ctx, waitGroup, config, store, redisDistro)

	if err := waitGroup.Wait(); err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
}

func runTaskProcessor(ctx context.Context, wg *errgroup.Group, redisOpt asynq.RedisClientOpt, store db.Store, config utils.Config)  {

	mailSender := mail.NewGmailSender(config.EmailSendName, config.From_Email, config.EamilPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailSender)

	slog.Info(
		"task taskProcessor started",
	)
	if err := taskProcessor.Start(); err != nil {
		slog.Error("failed to start task processor %w", 
		slog.String("error",err.Error()))
	}

	wg.Go(func() error {
		<-ctx.Done()
		slog.Info("shutting down processor")
		taskProcessor.Shutdown()

		return nil
	})
	
	
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

func startGinServer(ctx context.Context, wg *errgroup.Group, config utils.Config,store db.Store, taskDistributer worker.TaskDistributer)  {
	server, err := api.NewServer(config, store, taskDistributer)
	if err != nil {
		log.Fatal(err)
	}

	wg.Go(func() error {
		slog.Info("started the server", slog.String("port", config.HTTPServerAddress))
		err = server.StartServer(config.HTTPServerAddress); 
		if err != nil {
			slog.Error("error starting server", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<- ctx.Done()
		slog.Info("Gracefull shutdown")
		if err := server.ShutdownServer(context.Background()); err != nil {
			slog.Error("problem shutting down server", slog.String("error", err.Error()))
			return err
		}
		return nil

	})
	
}