package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phpunch/route-roam-api/log"
	"net/http"
	"time"
)

type postController interface {
	CreatePost(ctx *gin.Context)
	LikePost(ctx *gin.Context)
	UnlikePost(ctx *gin.Context)
	CommentPost(ctx *gin.Context)
}

func (c *controller) CreatePost(ctx *gin.Context) {
	userID, found := ctx.GetPostForm("userId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	var textPtr *string

	// upload images
	form, _ := ctx.MultipartForm()
	files := form.File["datasetPath[]"]

	filePathMinio := make([]string, 5)
	i := 0
	for _, file := range files {
		log.Log.Infof("upload file path: %s", file.Filename)
		objectName := userID + "/" + time.Now().Format("20060102") + "/" + file.Filename
		filepath, err := c.service.UploadFile(ctx, objectName, file, "image")
		if err != nil {
			ctx.JSON(http.StatusForbidden, fmt.Errorf("Failed upload file"))
		}
		filePathMinio[i] = filepath
		i = i + 1
	}

	// save metadata
	if err := c.service.CreatePost(userID, textPtr, filePathMinio); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"filepath": filePathMinio,
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

func (c *controller) UnlikePost(ctx *gin.Context) {
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
	if err := c.service.UnlikePost(userID, postID); err != nil {
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
