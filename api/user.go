package api

import (
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

type userCreateResponse struct {
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

	userResponse := userCreateResponse {
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, userResponse)




}