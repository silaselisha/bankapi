package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/silaselisha/bankapi/token"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationType       = "bearer"
	AuthorizationPayloadKey = "authorizationPayloadKey"
)

func AuthorizationMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("unauthorized invalid auth header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("unauthorized invalid auth header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fmt.Println(strings.ToLower(fields[0]))
		if strings.ToLower(fields[0]) != AuthorizationType {
			err := errors.New("unauthorized invalid auth type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		token := fields[1]
		payload, err := maker.VerifyToken(token)
		if err != nil {
			fmt.Println(token)
			err := errors.New("unauthorized invalid access token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
