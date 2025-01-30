package api

import (
	"database/sql"
	"net/http"

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
)

type createAcountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD KES"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var req createAcountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	
	}

	account, err := server.store.CreateAccount(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}


type getAccountRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}


func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {

		if err == sql.ErrNoRows {

		ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		
	}
	
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID int32   `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	args := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := server.store.ListAccounts(ctx, args)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		
	}
	
	ctx.JSON(http.StatusOK, account)
}



