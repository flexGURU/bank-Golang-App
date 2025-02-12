package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
)

type tokenRenewRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type tokenRenewResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}



func (server *Server) renewAccessToken(ctx *gin.Context) {

	var req tokenRenewRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	resfreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken) 
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	session, err  := server.store.GetSession(ctx, resfreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session user ")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	if session.Username != resfreshPayload.Username {
		err := fmt.Errorf("uknown session user")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("uknown session user")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}


	access_token, access_token_payload, err := server.tokenMaker.CreateToken(resfreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	Response := tokenRenewResponse {
		AccessToken: access_token,
		AccessTokenExpiresAt: access_token_payload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, Response)

}
