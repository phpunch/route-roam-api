package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phpunch/route-roam-api/log"
	"net/http"
	"time"
)

type fileController interface {
	UploadFiles(c *gin.Context)
	GetFile(ctx *gin.Context)
}

func (c *controller) UploadFiles(ctx *gin.Context) {
	dsID := "images"
	form, _ := ctx.MultipartForm()
	files := form.File["datasetPath[]"]

	filePathMinio := make([]string, 5)
	i := 0
	for _, file := range files {
		log.Log.Infof("upload file path: %s", file.Filename)
		objectName := time.Now().Format("20060102") + "/" + dsID + "/" + file.Filename
		filepath, err := c.service.UploadFile(ctx, objectName, file, "image")
		if err != nil {
			ctx.JSON(http.StatusForbidden, fmt.Errorf("Failed upload file"))
		}
		filePathMinio[i] = filepath
		i = i + 1
	}
	ctx.JSON(http.StatusOK, gin.H{
		"filepath": filePathMinio,
	})
}

func (c *controller) GetFile(ctx *gin.Context) {
	name := ctx.Query("name")
	dsID := "images"
	objectName := time.Now().Format("20060102") + "/" + dsID + "/" + name
	object, err := c.service.GetFile(ctx, objectName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	st, err := object.Stat()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.DataFromReader(http.StatusOK, st.Size, st.ContentType, object, nil)
}
