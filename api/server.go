package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	db "github.com/flexGURU/simplebank/db/sqlc"
	docs "github.com/flexGURU/simplebank/docs"
	"github.com/flexGURU/simplebank/token"
	"github.com/flexGURU/simplebank/utils"
	"github.com/flexGURU/simplebank/worker"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"

)



type Server struct {
	config utils.Config
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
	taskDistributer worker.TaskDistributer
	httpServer *http.Server
}

// NewServer will create a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store, taskDistributer worker.TaskDistributer) (*Server, error)  {
	docs.SwaggerInfo.BasePath = ""


	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) 
	if err != nil {
		return nil, fmt.Errorf("cannot create a token: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
		taskDistributer: taskDistributer,
	}

	server.serverRoutes()


	return server, nil
	
}

func (server *Server) serverRoutes() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{server.config.Origin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/renewtoken", server.renewAccessToken)
	router.GET("/verify_email", server.verifyEmail)


	authRoutes := router.Group("/").Use( authMiddleware(server.tokenMaker))

	// Adding the routes
	authRoutes.POST("/createaccount", server.createAccount)
	authRoutes.POST("/getaccount/:id", server.getAccount)
	authRoutes.GET("/listaccounts", server.listAccounts)
	authRoutes.GET("/users", server.listUsers)
	
	authRoutes.POST("/transfers", server.createTransfer)


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	

	server.router = router


}

// starting the server
func (server *Server) StartServer(address string) error {

	server.httpServer = &http.Server{
		Addr: address,
		Handler: server.router,
	}

	slog.Info(
		"starting http server",
		slog.String("address", address),
	)

	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed{
		slog.Error("failed to start Server", 
		slog.String("error starting", err.Error()))
		return err
	}
	return nil
	
}

func (server *Server) ShutdownServer(ctx context.Context) error {

	return server.httpServer.Shutdown(ctx);
	
}