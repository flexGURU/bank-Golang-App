package api

import (
	"fmt"

	db "github.com/flexGURU/simplebank/db/sqlc"
	docs "github.com/flexGURU/simplebank/docs"
	"github.com/flexGURU/simplebank/token"
	"github.com/flexGURU/simplebank/utils"
	"github.com/flexGURU/simplebank/worker"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)



type Server struct {
	config utils.Config
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
	taskDistributer worker.TaskDistributer
}

// NewServer will create a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store, taskDistributer worker.TaskDistributer) (*Server, error)  {
	docs.SwaggerInfo.BasePath = ""


	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) 
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

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/renewtoken", server.renewAccessToken)
	router.GET("/verify_email", server.verifyEmail)


	authRoutes := router.Group("/").Use( authMiddleware(server.tokenMaker))

	// Adding the routes
	authRoutes.POST("/createaccount", server.createAccount)
	authRoutes.POST("/getaccount/:id", server.getAccount)
	authRoutes.GET("/listaccounts", server.listAccounts)
	
	authRoutes.POST("/transfers", server.createTransfer)


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	

	server.router = router


}

// starting the server
func (server *Server) StartServer(address string) error {

	return server.router.Run(address)
	
}