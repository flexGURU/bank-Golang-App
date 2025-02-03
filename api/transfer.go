package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/flexGURU/simplebank/db/sqlc"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
)

type createTransferParams struct {
	FromAccountID int32  `json:"from_account_id" binding:"required,min=0"`
	ToAccountID   int32  `json:"to_account_id" binding:"required,min=0"`
	Amount        int32  `json:"amount" binding:"required,min=0"`
	Currency      string `json:"currency" binding:"required,oneof=KES USD"`
}

func (server *Server) createTransfer(ctx *gin.Context) {

	var req createTransferParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	args := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId: req.ToAccountID,
		Amount: req.Amount,
		
	}

	account, err := server.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

func (server Server) validAccount(ctx *gin.Context, accountID int32, currency string) bool  {

	
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {

		if err == sql.ErrNoRows {

		ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
		return false

		}

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return false
		
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mistmatch: %s vs  %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return false
	}

	return true
	
}
