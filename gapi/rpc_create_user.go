package gapi

import (
	"context"

	"github.com/flexGURU/simplebank/auth"
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	
	password, err := auth.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s, ", err)
	}

	args := db.CreateUserParams {
		Username: req.GetUsername(),
		HashedPassword: password,
		FullName: req.GetFullName(),
		Email: req.GetEmail(),
	}

	user, err  := server.store.CreateUser(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	userResponse := &pb.CreateUserResponse {
		User: convertUser(user),
	}

	return userResponse, nil
}