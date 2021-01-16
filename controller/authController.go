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
	Refresh(ctx *gin.Context)
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
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "email not found",
		})
		return
	}
	password, found := ctx.GetPostForm("password")
	if !found {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "password not found",
		})
		return
	}

	userID, err := c.service.LoginUser(email, password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
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

func (c *controller) LogoutUser(ctx *gin.Context) {
	accessUUID := ctx.GetString("access_uuid")
	deleted, delErr := c.service.DeleteAuth(accessUUID)
	if delErr != nil || deleted == 0 {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	ctx.JSON(http.StatusOK, "Successfully logged out")
}

func (c *controller) Refresh(ctx *gin.Context) {
	mapToken := map[string]string{}
	if err := ctx.ShouldBindJSON(&mapToken); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	token, err := c.service.VerifyToken(refreshToken, "REFRESH_SECRET")
	//if there is an error, the token must have expired
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
		return
	}
	tokens, err := c.service.RefreshToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
	}
	ctx.JSON(http.StatusCreated, tokens)
}
