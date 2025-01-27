package api

import (
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store *db.Store
	router *gin.Engine
}

// NewServer will create a new HTTP server and setup routing
func NewServer(store *db.Store) *Server  {
	server := &Server{store: store}
	router := gin.Default()

	// Adding the routes
	router.POST("/createaccount", createAcount)

	server.router = router

	return server
	
}