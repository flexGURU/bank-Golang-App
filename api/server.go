package api

import (

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)



type Server struct {
	store db.Store
	router *gin.Engine
}

// NewServer will create a new HTTP server and setup routing
func NewServer(store db.Store) *Server  {
	server := &Server{store: store}
	router := gin.Default()

	// Adding the routes
	router.POST("/createaccount", server.createAccount)
	router.POST("/getaccount/:id", server.getAccount)
	router.GET("/listaccounts", server.listAccounts)
	

	server.router = router
	

	return server
	
}

// starting the server
func (server *Server) StartServer(address string) error {

	return server.router.Run(address)
	
}