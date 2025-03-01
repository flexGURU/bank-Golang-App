package api

import (
	"net/http"

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
)

type verifyEmailRequest struct {
	Id int32 `form:"id" binding:"required,min=1"`
	SecretCode string `form:"code" binding:"required,min=1"`
}



func (server *Server) verifyEmail(ctx *gin.Context) {

	var req verifyEmailRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	args := db.VerifyEmailTxParams{
		UpdateVerifyEmailParams: db.UpdateVerifyEmailParams{
			ID: req.Id,
			Secretcode: req.SecretCode,
		},
	}

	result, err := server.store.VerifyEmailTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	response := result.IsVerified

	ctx.JSON(http.StatusAccepted, response)





}