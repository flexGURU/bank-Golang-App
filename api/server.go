package api

import (
	"fmt"

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/token"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
   	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/flexGURU/simplebank/docs"

)



type Server struct {
	config utils.Config
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
}

// NewServer will create a new HTTP server and setup routing
func NewServer(config utils.Config,store db.Store) (*Server, error)  {
	docs.SwaggerInfo.BasePath = ""


	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) 
	if err != nil {
		return nil, fmt.Errorf("cannot create a token: %w", err)
	}
	
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	server.serverRoutes()

	return server, nil
	
}

func (server *Server) serverRoutes() {
	router := gin.Default()

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/renewtoken", server.renewAccessToken)

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