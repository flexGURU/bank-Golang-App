package gapi

import (
	"context"
	"database/sql"

	"github.com/flexGURU/simplebank/auth"
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)



func (server *Server) LoginUser(ctx context.Context,req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	
	
	user, err  := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, status.Errorf(codes.NotFound, "user not found")
			
		}
		return nil, status.Errorf(codes.Internal, "Internal Server Error")

	}

	if err := auth.ComparePassword(user.HashedPassword, req.Password); err != nil {
		return nil, status.Errorf(codes.NotFound, "Invalid Password")

	}

	access_token, access_token_payload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create access token")
	}

	referesh_token, referesh_token_payload,  err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create access token")
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams {
		ID: referesh_token_payload.ID,
		Username: referesh_token_payload.Username,
		RefreshToken: referesh_token,
		UserAgent: "",
		ClientIp: "",
		IsBlocked: false,
		ExpiresAt: referesh_token_payload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session")
	}

	response := &pb.LoginUserResponse{
		User: convertUser(user),
		SessionID: session.ID.Version().String(),
		AccessToken: access_token,
		RefreshToken: referesh_token,
		AccessTokenExpiresAt: timestamppb.New(access_token_payload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(referesh_token_payload.ExpiredAt),
	}

	return response, nil
}