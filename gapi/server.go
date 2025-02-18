package gapi

import (
	"fmt"

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/pb"
	"github.com/flexGURU/simplebank/token"
	"github.com/flexGURU/simplebank/utils"
)


type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer will create a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create a token: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}


	return server, nil

}