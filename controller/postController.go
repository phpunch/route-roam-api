package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type postController interface {
	CreatePost(ctx *gin.Context)
	LikePost(ctx *gin.Context)
	CommentPost(ctx *gin.Context)
}

func (c *controller) CreatePost(ctx *gin.Context) {
	userID, found := ctx.GetPostForm("userId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	var textPtr *string
	var imagePtr *string
	*textPtr, _ = ctx.GetPostForm("text")
	*imagePtr, _ = ctx.GetPostForm("image")

	if err := c.service.CreatePost(userID, textPtr, imagePtr); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}
func (c *controller) LikePost(ctx *gin.Context) {
	userID, found := ctx.GetPostForm("userId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	postID, found := ctx.GetPostForm("postId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	if err := c.service.LikePost(userID, postID); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func (c *controller) CommentPost(ctx *gin.Context) {
	userID, found := ctx.GetPostForm("userId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	postID, found := ctx.GetPostForm("postId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	text, found := ctx.GetPostForm("text")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	if err := c.service.CommentPost(userID, postID, text); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}
