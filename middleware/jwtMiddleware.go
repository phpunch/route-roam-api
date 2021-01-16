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
		ad, err := m.service.ExtractTokenMetadata(ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": fmt.Sprintf("unauthorized: %v", err),
			})
			return
		}
		userID, err := m.service.FetchAuth(ad)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": fmt.Sprintf("unauthorized: %v", err),
			})
			return
		}
		ctx.Set("user_id", userID)
		ctx.Set("access_uuid", ad.AccessUUID)
	}
}
