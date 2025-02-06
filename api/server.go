package api

import (
	"fmt"

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/token"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
)



type Server struct {
	config utils.Config
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
}

// NewServer will create a new HTTP server and setup routing
func NewServer(config utils.Config,store db.Store) (*Server, error)  {

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

	// Adding the routes
	router.POST("/createaccount", server.createAccount)
	router.POST("/getaccount/:id", server.getAccount)
	router.GET("/listaccounts", server.listAccounts)
	router.POST("/transfers", server.createTransfer)
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	server.router = router


}

// starting the server
func (server *Server) StartServer(address string) error {

	return server.router.Run(address)
	
}