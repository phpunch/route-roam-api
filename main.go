package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phpunch/route-roam-api/controller"
	"github.com/phpunch/route-roam-api/infrastructure/db"
	"github.com/phpunch/route-roam-api/log"
	"github.com/phpunch/route-roam-api/middleware"
	"github.com/phpunch/route-roam-api/repository"
	"github.com/phpunch/route-roam-api/service"
	"github.com/spf13/viper"
	"os"
)

func getLogger() log.Logger {
	logLevel := viper.GetString("Log.Level")
	logLevel = log.NormalizeLogLevel(logLevel)

	logColor := viper.GetBool("Log.Color")
	logJSON := viper.GetBool("Log.JSON")

	logger, err := log.NewLogger(&log.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logLevel,
		ConsoleJSONFormat: logJSON,
		Color:             logColor,
	}, log.InstanceZapLogger)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return logger
}

func main() {
	////////////// GIN //////
	logger := getLogger()
	dbHandler, err := db.NewDB(logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := repository.NewRepository(dbHandler)
	s := service.NewService(r)
	c := controller.NewController(s)
	mw := middleware.New(s)

	router := gin.Default()

	router.Use(middleware.CorMiddleware("*"))

	router.POST("/api/register", c.RegisterUser)
	router.POST("/api/login", c.LoginUser)
	router.POST("/api/token/refresh", c.Refresh)
	router.GET("/api/file/*filepath", c.GetFile)

	router.Use(mw.AuthorizeToken())
	{
		router.POST("/api/logout", c.LogoutUser)
		router.POST("/api/file", c.UploadFiles)
		router.POST("/api/post", c.CreatePost)
		router.DELETE("/api/post/:postId", c.DeletePost)
		router.POST("/api/like", c.LikePost)
		router.POST("/api/unlike", c.UnlikePost)
		router.GET("/api/posts", c.GetPosts)
		router.POST("/api/comment", c.CommentPost)
		router.GET("/api/comment/:postId", c.GetCommentsByPostID)
	}

	router.Run()
}
