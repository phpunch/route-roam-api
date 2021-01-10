package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	password, found := ctx.GetPostForm("password")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
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

	//verify the token
	// token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
	// 	//Make sure that the token method conform to "SigningMethodHMAC"
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return []byte(os.Getenv("REFRESH_SECRET")), nil
	// })
	token, err := c.service.VerifyToken2(refreshToken, "REFRESH_SECRET")
	//if there is an error, the token must have expired
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			ctx.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := c.service.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			ctx.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := c.service.CreateToken(userId)
		if createErr != nil {
			ctx.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := c.service.CreateAuth(userId, ts)
		if saveErr != nil {
			ctx.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		ctx.JSON(http.StatusCreated, tokens)
	} else {
		ctx.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
