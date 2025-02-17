package gapi

import (
	"context"

	"github.com/flexGURU/simplebank/pb"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}