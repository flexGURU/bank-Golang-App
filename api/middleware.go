package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/flexGURU/simplebank/token"
	"github.com/flexGURU/simplebank/utils"
	"github.com/gin-gonic/gin"
)


const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "bearer"
	authPayloadKey = "auth_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func (ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("no authorisation Header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("inavlid auth header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationTypeBearer {
			err := errors.New("inavlid auth type not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()

	}
}
