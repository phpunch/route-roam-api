package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authController interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
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
