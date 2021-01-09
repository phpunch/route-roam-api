package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phpunch/route-roam-api/service"
	"net/http"
)

type middleware struct {
	service service.Service
}

func New(service service.Service) *middleware {
	return &middleware{
		service,
	}
}

func (m *middleware) AuthorizeToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenAuth, err := m.service.ExtractTokenMetadata(ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Sprintf("unauthorized: %v", err))
			return
		}
		userID, err := m.service.FetchAuth(tokenAuth)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Sprintf("unauthorized: %v", err))
			return
		}
		ctx.Set("user_id", userID)
	}
}
