package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phpunch/route-roam-api/log"
	"net/http"
	"strconv"
	"time"
)

type postController interface {
	CreatePost(ctx *gin.Context)
	LikePost(ctx *gin.Context)
	UnlikePost(ctx *gin.Context)
	CommentPost(ctx *gin.Context)
	GetPosts(ctx *gin.Context)
}

func (c *controller) CreatePost(ctx *gin.Context) {
	tokenAuth, err := c.service.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, fmt.Sprintf("unauthorized: %v", err))
		return
	}
	userID, err := c.service.FetchAuth(tokenAuth)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, fmt.Sprintf("unauthorized: %v", err))
		return
	}

	// userIDStr, found := ctx.GetPostForm("userId")
	// if !found {
	// 	ctx.Status(http.StatusUnprocessableEntity)
	// 	return
	// }
	var textPtr *string

	// userID, err := strconv.Atoi(userIDStr)
	// if err != nil {
	// 	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"message": "userID is not number",
	// 	})
	// 	return
	// }

	// upload images
	form, _ := ctx.MultipartForm()
	files := form.File["datasetPath[]"]

	filePathMinio := make([]string, 5)
	i := 0
	for _, file := range files {
		log.Log.Infof("upload file path: %s", file.Filename)
		objectName := strconv.Itoa(int(userID)) + "/" + time.Now().Format("20060102") + "/" + file.Filename
		log.Log.Infof("ovjectname: %s", objectName)
		filepath, err := c.service.UploadFile(ctx, objectName, file, "image")
		if err != nil {
			ctx.JSON(http.StatusForbidden, fmt.Errorf("Failed upload file"))
		}
		filePathMinio[i] = filepath
		i = i + 1
	}

	// save metadata
	if err := c.service.CreatePost(int(userID), textPtr, filePathMinio); err != nil {
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
	userIDStr, found := ctx.GetPostForm("userId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	postIDStr, found := ctx.GetPostForm("postId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "userID is not number",
		})
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "userID is not number",
		})
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
	userIDStr, found := ctx.GetPostForm("userId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	postIDStr, found := ctx.GetPostForm("postId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "userID is not number",
		})
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "userID is not number",
		})
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
	userIDStr, found := ctx.GetPostForm("userId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	postIDStr, found := ctx.GetPostForm("postId")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}
	text, found := ctx.GetPostForm("text")
	if !found {
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "userID is not number",
		})
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "userID is not number",
		})
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

func (c *controller) GetPosts(ctx *gin.Context) {
	posts, err := c.service.GetPosts()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"posts":   posts,
	})
}
