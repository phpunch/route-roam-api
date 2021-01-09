package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authController interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	LogoutUser(ctx *gin.Context)
}

func (c *controller) RegisterUser(ctx *gin.Context) {
	email, found := ctx.GetPostForm("email")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	password, found := ctx.GetPostForm("password")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	userID, err := c.service.RegisterUser(email, password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	token, err := c.service.CreateToken(userID)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = c.service.CreateAuth(userID, token)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func (c *controller) LoginUser(ctx *gin.Context) {
	email, found := ctx.GetPostForm("email")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	password, found := ctx.GetPostForm("password")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	if err := c.service.LoginUser(email, password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func (c *controller) LogoutUser(ctx *gin.Context) {
	accessUUID := ctx.GetString("access_uuid")
	deleted, delErr := c.service.DeleteAuth(accessUUID)
	if delErr != nil || deleted == 0 {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	ctx.JSON(http.StatusOK, "Successfully logged out")
}
