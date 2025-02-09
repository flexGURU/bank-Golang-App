package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/flexGURU/simplebank/auth"
	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
)


type createUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	HashedPassword string `json:"hashed_password" binding:"required,min=6"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	password, err := auth.HashPassword(req.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	args := db.CreateUserParams {
		Username: req.Username,
		HashedPassword: password,
		FullName: req.FullName,
		Email: req.Email,
	}

	user, err  := server.store.CreateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	userResponse := userResponse {
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, userResponse)

}


type loginUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	HashedPassword string `json:"hashed_password" binding:"required,min=6"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	User userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {

	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err  := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if err := auth.ComparePassword(user.HashedPassword, req.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	access_token, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
	}

	Response := loginResponse {
		AccessToken: access_token,
		User: userResponse{
			Username: user.Username,
			FullName: user.FullName,
			Email: user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt: user.CreatedAt,
			},
		
	}

	ctx.JSON(http.StatusOK, Response)

}
