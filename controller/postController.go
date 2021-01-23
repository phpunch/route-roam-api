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
	GetCommentsByPostID(ctx *gin.Context)
}

func (c *controller) CreatePost(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	text, found := ctx.GetPostForm("text")
	if !found {
		ctx.JSON(http.StatusUnprocessableEntity, "text is not found")
		return
	}

	// upload images
	form, _ := ctx.MultipartForm()
	files := form.File["datasetPath[]"]

	filePathMinio := []string{}
	for _, file := range files {
		log.Log.Infof("upload file path: %s", file.Filename)
		objectName := strconv.Itoa(int(userID)) + "/" + time.Now().Format("20060102") + "/" + file.Filename
		log.Log.Infof("ovjectname: %s", objectName)
		filepath, err := c.service.UploadFile(ctx, objectName, file, "image")
		if err != nil {
			ctx.JSON(http.StatusForbidden, fmt.Errorf("Failed upload file"))
		}
		filePathMinio = append(filePathMinio, filepath)
	}

	// save metadata
	post, err := c.service.CreatePost(userID, text, filePathMinio)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"post":    post,
	})

}
func (c *controller) LikePost(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")

	postIDStr, found := ctx.GetPostForm("postId")
	if !found {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "postId is not found",
		})
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "postId is not number",
		})
		return
	}

	if err := c.service.LikePost(userID, int64(postID)); err != nil {
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
	userID := ctx.GetInt64("user_id")

	postIDStr, found := ctx.GetPostForm("postId")
	if !found {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "postId is not found",
		})
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "postId is not number",
		})
		return
	}
	if err := c.service.UnlikePost(userID, int64(postID)); err != nil {
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
	userID := ctx.GetInt64("user_id")

	postIDStr, found := ctx.GetPostForm("postId")
	if !found {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "postId is not found",
		})
		return
	}
	text, found := ctx.GetPostForm("text")
	if !found {
		ctx.JSON(http.StatusUnprocessableEntity, "text is not found")
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "userID is not number",
		})
		return
	}
	comment, err := c.service.CommentPost(userID, int64(postID), text)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"comment": comment,
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

func (c *controller) GetCommentsByPostID(ctx *gin.Context) {
	postIDStr := ctx.Param("postId")
	fmt.Printf("postId: %v\n", postIDStr)
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "postId is not number",
		})
		return
	}

	comments, err := c.service.GetCommentsByPostID(int64(postID))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"comments": comments,
	})
}
